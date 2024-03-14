package main

import (
	"context"
	"encoding/csv"
	"encoding/json"
	"fmt"
	"os"
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
		Addrs: []string{"localhost:6376", "localhost:6375"},
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
	defer func() {
		writer.Flush()
	}()

	var cursor uint64
	var n int
	for {
		var keys []string
		keys, cursor, err = rdb.Scan(ctx, cursor, "*", 0).Result()
		if err != nil {
			fmt.Println("error scanning keys:", err)
			return
		}

		for _, key := range keys {
			val, err := rdb.Get(ctx, key).Result()
			if err != nil {
				fmt.Println("error getting value for key:", key, err)
				continue
			}

			var data Data
			err = json.Unmarshal([]byte(val), &data)
			if err != nil {
				fmt.Println("error unmarshalling json:", err)
				continue
			}

			if data.AssignmentDate == yesterday && data.IsExcluded {
				if n == 0 {
					headers := []string{"assignmentDate", "isExcluded"}
					if err := writer.Write(headers); err != nil {
						fmt.Println("nth error:", err)
					}
				}

				record := []string{data.AssignmentDate, fmt.Sprintf("%t", data.IsExcluded)}
				if err := writer.Write(record); err != nil {
					fmt.Println("record error:", err)
				}
				n++
			}
		}

		if cursor == 0 {
			break
		}
	}

	fmt.Printf("csv file %s created, %d records written.\n", fileName, n)
}
