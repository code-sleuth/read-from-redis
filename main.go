package main

import (
	"context"
	"encoding/csv"
	"encoding/json"
	"fmt"
	"os"
	"sync"
	"time"

	"github.com/redis/go-redis/v9"
)

type Data struct {
	AssignmentDate string `json:"assignmentDate"`
	IsExcluded     bool   `json:"isExcluded"`
}

func main() {
	ctx := context.Background()
	rdb := redis.NewClusterClient(&redis.ClusterOptions{
		Addrs: []string{"localhost:7002", "localhost:7003"},
	})

	_, err := rdb.Ping(ctx).Result()
	if err != nil {
		fmt.Println("error connecting to redis cluster:", err)
		return
	}

	// handle csv file
	yesterday := time.Now().AddDate(0, 0, -1).Format("2006-01-02")
	fileName := fmt.Sprintf("%s_results.csv", time.Now().Format("200601021504"))
	file, err := os.Create(fileName)
	if err != nil {
		fmt.Println("error creating csv file:", err)
		return
	}
	defer file.Close()

	writer := csv.NewWriter(file)
	defer writer.Flush()

	var wg sync.WaitGroup
	dataChannel := make(chan Data)
	errChannel := make(chan error, 1)
	doneChannel := make(chan bool)

	go func() {
		var n int
		for data := range dataChannel {
			if n == 0 {
				headers := []string{"assignmentDate", "isExcluded"}
				if err := writer.Write(headers); err != nil {
					errChannel <- err
					return
				}
			}
			record := []string{data.AssignmentDate, fmt.Sprintf("%t", data.IsExcluded)}
			if err := writer.Write(record); err != nil {
				errChannel <- err
				return
			}
			n++
		}
		doneChannel <- true
	}()

	var cursor uint64
	cap := make(chan struct{}, 50)
	for {
		var keys []string
		keys, cursor, err = rdb.Scan(ctx, cursor, "*", 0).Result()
		if err != nil {
			fmt.Println("error scanning keys:", err)
			return
		}

		for _, key := range keys {
			wg.Add(1)
			cap <- struct{}{}
			go func(key string) {
				defer wg.Done()
				val, err := rdb.Get(ctx, key).Result()
				if err != nil {
					fmt.Println("error getting value for key:", key, err)
					<-cap
					return
				}

				var data Data
				if err := json.Unmarshal([]byte(val), &data); err != nil {
					fmt.Println("error unmarshalling json:", err)
					<-cap
					return
				}

				if data.AssignmentDate == yesterday && data.IsExcluded {
					dataChannel <- data
				}
				<-cap
			}(key)
		}

		if cursor == 0 {
			break
		}
	}

	wg.Wait()
	close(dataChannel)
	select {
	case <-doneChannel:
		fmt.Printf("csv file %s created.\n", fileName)
	case err := <-errChannel:
		fmt.Printf("error writing to csv: %v\n", err)
	}
}
