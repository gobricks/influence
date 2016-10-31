package influence

import (
	"net/http"
	"time"

	client "github.com/influxdata/influxdb/client/v2"
)

// DefaultMiddleware sends http handler statistic to InfluxDB
func DefaultMiddleware(conn client.Client, next http.Handler) http.Handler {
	tags := getEnvTags(confHandlerTagsKey)

	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
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

			tags["handler"] = r.URL.Path

			pt, err := client.NewPoint(httpHandlerTablename, tags, fields, endTime)
			if err != nil {
				return
			}

			bp.AddPoint(pt)
			conn.Write(bp)
		}()

		next.ServeHTTP(w, r)
	})
}
