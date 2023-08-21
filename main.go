package main

import (
	"cribl/logmonitor/file"
	"github.com/gin-gonic/gin"
	"net/http"
	"strconv"
)

const FILE_PATH = "/var/log/"

func main() {
	router := gin.Default()

	router.GET("/api/v1/logs", func(c *gin.Context) {
		size := c.DefaultQuery("size", "100")
		filename := c.DefaultQuery("filename", "var5MB.txt")
		searchKeyword := c.Query("keyword")

		numOfEntries, err := strconv.Atoi(size)
		if err != nil {
			c.Error(err)
		}

		filenameWithPath := FILE_PATH + filename

		result, err := file.ReadLastNLinesWithKeyword(filenameWithPath, numOfEntries, searchKeyword)
		if err != nil {
			c.Error(err)
		}

		c.Writer.Header().Set("Access-Control-Allow-Origin", "*")
		c.IndentedJSON(http.StatusOK, result)
	})

	router.GET("/api/v1/plogs", func(c *gin.Context) {
		size := c.DefaultQuery("size", "100")
		filename := c.DefaultQuery("filename", "var5MB.txt")
		searchKeyword := c.Query("keyword")

		numOfEntries, err := strconv.Atoi(size)
		if err != nil {
			c.Error(err)
		}

		filenameWithPath := FILE_PATH + filename

		result, err := file.ReadLastNLinesWithKeywordP(filenameWithPath, numOfEntries, searchKeyword)
		if err != nil {
			c.Error(err)
		}

		retVal := make([]string, len(result))
		for i, r := range result {
			// strip line break at the end so that response is human readable
			retVal[i] = string(r[:len(r)-1])
		}
		c.Writer.Header().Set("Access-Control-Allow-Origin", "*")
		c.IndentedJSON(http.StatusOK, retVal)
	})

	router.Run("localhost:8080")
}
