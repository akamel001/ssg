package redis
import (
	"log"
	"strings"
	"menteslibres.net/gosexy/redis"
	"fmt"
)

var client *redis.Client

func Ssg_redis() {
	var err error

	var mmm map[string]string
	mmm = make(map[string]string)

	stats_read := []string{"CPU", "Memory", "Stats"}

	client = redis.New()

	err = client.Connect("localhost", 6379)

	if err != nil {
		log.Fatalf("Connection Failed: %s\n", err.Error())
		return
	}

	log.Println("Connected to Redis Server")

	log.Println("Getting INFO")

	// Random thought; should the metric to be monitored be rad from the config file; metric:file_path
	for key := range stats_read {
		info, err := client.Info(stats_read[key])

		if err !=nil {
			log.Fatalf("Failed to get data from Redis")

			// if info in not 0, start parsing
		}else {
			data_split(stats_read[key], info, mmm)

		}
	}
	for key, value := range mmm {
		fmt.Println("Key:", key, "Value:", value)
	}

}

func data_split(delimiter, info string, mmm map[string]string){
	delimiter = "#" + " " + delimiter
	result := strings.Split(info, delimiter)
	k, v := result[0], result[1]
	mmm[k]=v

}


