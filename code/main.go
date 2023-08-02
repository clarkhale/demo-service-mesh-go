package main

import (
	"encoding/json"
	"flag"
	"io/ioutil"
	"net/http"
	"os"
	"fmt"
	"strconv"

	"github.com/gin-gonic/gin"
	"github.com/penglongli/gin-metrics/ginmetrics"
)

var count int = 0
var hostname string
var nextHop *string = nil
var version *string = nil
var side *string = nil
var errCount = 0

type ExampleResponse struct {
	Count           int              `json:"count"`
	Hostname        string           `json:"hostname"`
	Side            string           `json:"side"`
	Version         string           `json:"version"`
	DownStreamError bool             `json:"downstreamerror"`
	NestedResponse  *ExampleResponse `json:"nestedResponse"`
}

func getCount(c *gin.Context) {
	count++

	response := ExampleResponse{Count: count,
		Hostname:        hostname,
		Side:            *side,
		Version:         *version,
		DownStreamError: false,
		NestedResponse:  nil}

	if *nextHop != "" {
		nextHopResponse, err := http.Get("http://" + *nextHop)
		fmt.Printf("Curling %s\n", *nextHop);
		defer nextHopResponse.Body.Close()

		if err != nil || nextHopResponse.Body == nil {
			errCount++
			response.DownStreamError = true
		} else {
			body, err := ioutil.ReadAll(nextHopResponse.Body)

			if err != nil {
				errCount++
				response.DownStreamError = true
			} else {
				var nestedResponse ExampleResponse
				err := json.Unmarshal(body, &nestedResponse)
				if err != nil {
					errCount++
					response.DownStreamError = true
				} else {
					response.NestedResponse = &nestedResponse
				}
			}
		}
	}
	c.IndentedJSON(http.StatusOK, response)
}

func main() {

	port := flag.Int("port", 8080, "Set Port Number")
	//nextHop = flag.String("nexthop", "localhost:8081", "Specify Next Hop")

	flag.Parse()

	bindHost := ":" + strconv.Itoa(*port)

	hostname = os.Getenv("HOSTNAME")

	nextHopEnv := os.Getenv("NEXT_HOP")

	nextHop = &nextHopEnv

	versionEnv := os.Getenv("APP_VERSION")

	version = &versionEnv

	sideEnv := os.Getenv("APP_SIDE")

	side = &sideEnv

	router := gin.Default()

	m := ginmetrics.GetMonitor()

	m.SetMetricPath("/metrics")

	m.Use(router)

	router.GET("/", getCount)

	router.Run(bindHost)
}
