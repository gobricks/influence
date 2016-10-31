package influence

import (
	"os"
	"strconv"
	"time"
)

var (
	confDatabase = os.Getenv("INFLUX_DATABASE")

	confHandlerTagsKey = os.Getenv("INFLUX_HTTP_HANDLER_TAGS")
	confRuntimeTagsKey = os.Getenv("INFLUX_GO_RUNTIME_TAGS")

	monitorAggregationMilliseconds, _ = strconv.Atoi(get(os.Getenv("INFLUX_AGGREGATION_INTERVAL"), "100"))
	confMonitoringAggregationInterval = time.Duration(monitorAggregationMilliseconds) * time.Millisecond

	monitorSyncMilliseconds, _ = strconv.Atoi(get(os.Getenv("INFLUX_SYNC_INTERVAL"), "1000"))
	confMonitoringSyncInterval = time.Duration(monitorSyncMilliseconds) * time.Millisecond
)

const (
	httpHandlerTablename = "http_response"
	runtimeTablename     = "go_runtime"
	gcTablename          = "go_gc"
)

// get is a ternary operator for string values.
// It returns first param in it is not empty string, else - returns second param.
func get(check string, dflt string) string {
	if check == "" {
		return dflt
	}
	return check
}
