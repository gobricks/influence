# Influence
Influence is a convenient package to reproduce New Relic go-agent functionality with Influx database.

# Install

```
go get -u github.com/gobricks/influence
```

# Setup
First of all you need to set few environvent variables

**Required:**
* **$INFLUX_DATABASE** - name of database to write data into.
* **$INFLUX_AGGREGATION_INTERVAL** - time interval to gather runtime statistics in milliseconds. Default: _100_.
* **$INFLUX_SYNC_INTERVAL** - time interval to save gathered runtime statistics to Influx in milliseconds. Must be at least two times greater than `$INFLUX_AGGREGATION_INTERVAL` Default: _1000_.

**Optional:**
* **$INFLUX_HTTP_HANDLER_TAGS** - tags to send with each gathered http handler statistics record. Format: comma separated "key:value", e.g. `hostname:dev.local,dc:us-west1`.
* **$INFLUX_GO_RUNTIME_TAGS** - tags to send with each gathered runtime statistics record. Format: same as above.

# Usage

```go
import (
    "net/http"
    "github.com/gobricks/influence"
    client "github.com/influxdata/influxdb/client/v2"
    "github.com/gin-gonic/gin"
)

func main() {
    r := gin.New()
    
    // connect to InfluxDB
    influxConn, err = client.NewHTTPClient(client.HTTPConfig{
		Addr:     "http://localhost:8086",
		Username: "user",
		Password: "pass",
	})
    if err != nil {
        panic(err)
    }
    
    // setup influence middleware
    r.Use(influence.GinMiddleware(influxConn))
    r.GET("/hello", hello)
    
    // start runtime statistics aggregation
    influence.StartMonitoring(influxConn)
    
    r.Run("127.0.0.1:8000")
}

func hello(c *gin.Context) {
	c.String(http.StatusOK, "Hello, world!")
}
```

# Tables
This package will create three main tables in your Influx database:

**http_response** - holds data for http response execution time. Note that handler_func only available in gin-gonic middleware. Execution time measures in milliseconds.
```
time			    execution_time	handler	 handler_func
1477580561978794221	28		        /hello	 hello
```

**go_runtime** - holds data for runtime memory and goroutines statistics.
```
time			     cpu_count	goroutines_count  memory_alloc
1477580561876964641	 12		    10			      2169480
```

**go_gc** - holds data for runtime garbage collector pauses.
```
time			     duration_ns
1477913591906347048	 366395
```

# Visual example
[![GoDoc](http://leproimg.com/2541851)](http://leproimg.com/2541851)