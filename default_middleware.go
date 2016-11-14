package influence

import (
	"net/http"
	"strconv"
	"time"

	client "github.com/influxdata/influxdb/client/v2"
)

// DefaultMiddleware sends http handler statistic to InfluxDB
func DefaultMiddleware(conn client.Client, next http.Handler) http.Handler {
	tags := getEnvTags(confHandlerTagsKey)

	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		startTime := time.Now()

		defer func() {
			go func() {
				endTime := time.Now()

				bp, err := getPointsBatch()
				if err != nil {
					return
				}

				contentLength, err := strconv.Atoi(w.Header().Get("Content-Length"))
				if err != nil {
					contentLength = 0
				}

				fields := map[string]interface{}{
					"execution_time": endTime.Sub(startTime).Nanoseconds() / int64(time.Millisecond),
					"content_length": contentLength,
				}

				handlerTags := map[string]string{
					"handler": r.URL.Path,
					"method":  r.Method,
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

		next.ServeHTTP(w, r)
	})
}
