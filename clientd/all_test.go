package main
import (
	"flag"
	"fmt"
	"log"
	"testing"
)

var (
	testHost string
	testPort uint
)

func init() {
	// For now hard coding input

	testHost = "localhost"
	testPort = 6379

	log.Printf("Running tests")

}


func TestConnect(t *Testing.T) {
	var err error
	client = New()

	//Attempting a connection
	err = client.Connect(testHost, testPort)
	if err != nil {
		t.Fatalf("Failed to connect to test server")
	}
}

