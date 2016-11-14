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
			go func() {
				endTime := time.Now()

				bp, err := getPointsBatch()
				if err != nil {
					return
				}

				fields := map[string]interface{}{
					"execution_time": endTime.Sub(startTime).Nanoseconds() / int64(time.Millisecond),
					"content_length": c.Writer.Size(),
					"status_code":    c.Writer.Status(),
				}

				handlerTags := map[string]string{
					"handler":      c.Request.URL.Path,
					"handler_func": filepath.Base(c.HandlerName()),
					"method":       c.Request.Method,
				}

				for k, v := range handlerTags {
					tags[k] = v
				}

				pt, err := client.NewPoint(httpHandlerTablename, tags, fields, endTime)
				if err != nil {
					return
				}

				bp.AddPoint(pt)
				conn.Write(bp)
			}()
		}()

		c.Next()
	}
}
