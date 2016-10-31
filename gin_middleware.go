package influence

import (
	"path/filepath"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/influxdata/influxdb/client/v2"
)

// GinMiddleware sends http handler statistic to InfluxDB
func GinMiddleware(conn client.Client) gin.HandlerFunc {
	tags := getEnvTags(confHandlerTagsKey)

	return func(c *gin.Context) {
		startTime := time.Now()

		defer func() {
			endTime := time.Now()

			bp, err := getPointsBatch()
			if err != nil {
				return
			}

			fields := map[string]interface{}{
				"execution_time": endTime.Sub(startTime).Nanoseconds() / int64(time.Millisecond),
			}

			tags["handler"] = c.Request.URL.Path
			tags["handler_func"] = filepath.Base(c.HandlerName())

			pt, err := client.NewPoint(httpHandlerTablename, tags, fields, endTime)
			if err != nil {
				return
			}

			bp.AddPoint(pt)
			conn.Write(bp)
		}()

		c.Next()
	}
}
