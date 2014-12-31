package redis
import (
	"log"
	"strings"
	"fmt"
	"menteslibres.net/gosexy/redis"
)

var client *redis.Client

func Ssg_redis() {
	var err error
	client = redis.New()

	err = client.Connect("localhost", 6379)

	if err != nil {
		log.Fatalf("Connection Failed: %s\n", err.Error())
		return
	}

	log.Println("Connected to Redis Server")

	log.Println("Getting INFO")

	// Random thought; should the metric to be monitored be rad from the config file; metric:file_path
	info, err := client.Info("CPU")

	if len(info) == 0 {
		log.Fatalf("Failed to get data from Redis")

		// if info in not 0, start parsing
	}else {
		data_split("CPU", info)
	}

}

func data_split(delimiter string, info string) {
	delimiter = "#" + " " + delimiter
	result := strings.Split(info, delimiter)
	k, v := result[0], result[1]
	fmt.Println(k, v)

}

