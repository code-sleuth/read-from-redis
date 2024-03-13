# Read Keys From Redis With Cluster Mode Enabled

The following implementation creates a cluster with 3 master and 0 replicas.

### Optional
- **Docker Compose Configuration**:
    - Firstly, ensure your `docker-compose.yml` file contains a service block for Redis. This block defines how your Redis containers are set up, including their image, ports, and network settings. If such a block doesn't already exist, you'll need to create it, specifying the necessary configuration for each Redis node in your cluster.
- **Initialization Script Modification**:
	- Next, locate and edit the initialization script named `docker-entrypoint.sh` found in the `docker-data` directory. This script is executed when your Redis containers start, setting up the initial configuration of the Redis cluster.
	    
	- In this script, you'll find a line responsible for creating the Redis cluster using the `redis-cli --cluster create` command. The command looks like this:
    ```shell
    echo "yes" | redis-cli --cluster create 173.17.0.2:7002 173.17.0.3:7003 173.17.0.4:7004 --cluster-replicas 0
    ```
    - Additionally, change the `--cluster-replicas 0` option to `--cluster-replicas 1`. This modification alters the cluster's replica configuration from having no replicas (0) to having one replica per master node (1), enhancing your cluster's fault tolerance and data redundancy. You'll need to add 3 more IP addressess.

### Production
If you want to use this in production, define each Redis container to run on a static machine.

### How To Start
In your terminal
1. run `$ docker-compose build`
2. run`$ docker-compose up` or `$ docker-compose up -d`

## Test Script
In another terminal
1. run `$ go mod tidy`
2. run `$ go run main.go`
3. check the generated csv file for the redis keys
