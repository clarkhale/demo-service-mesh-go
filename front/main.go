package main

import (
        "net/http"
	"github.com/gin-gonic/gin"
	"github.com/penglongli/gin-metrics/ginmetrics"
)

var count int = 0

func getCount(c *gin.Context) {
	count++
        
        c.IndentedJSON(http.StatusOK, count)
}

func main() {
	router := gin.Default()

        m := ginmetrics.GetMonitor()

        m.SetMetricPath("/metrics")

        m.Use(router)

	router.GET("/", getCount)

	router.Run("localhost:8080")
}

