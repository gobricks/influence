package influence

import (
	"errors"
	"runtime"
	"runtime/debug"
	"time"

	client "github.com/influxdata/influxdb/client/v2"
)

// Set default last GC start to UNIX 0 time
var latestGCStart = time.Unix(0, 0)

// StartMonitoring spawns a goroutine and sends runtime service statistics into InfluxDB.
// This method will spawn two goroutines: one for runtime statistics aggregation
// (every *aggregateInterval* period) and one for syncing this data with InfluxDB
// (every *syncInterval* period)
func StartMonitoring(conn client.Client) error {
	if confMonitoringSyncInterval*2 < confMonitoringAggregationInterval {
		return errors.New("Sync interval must at least x2 longer than aggregation interval")
	}

	pointsChan := make(chan *client.Point)
	go dataFlusher(conn, pointsChan)
	go dataGatherer(pointsChan)

	return nil
}

// Syncs this data with InfluxDB (every *syncInterval* period)
func dataFlusher(conn client.Client, pointsChan chan *client.Point) {
	bp, _ := getPointsBatch()
	flushTicker := time.NewTicker(confMonitoringSyncInterval)

	for {
		select {
		case pt := <-pointsChan:
			// store received point
			bp.AddPoint(pt)
		case <-flushTicker.C:
			// clone gathered points
			sbp, _ := getPointsBatch()
			sbp.AddPoints(bp.Points())

			// renew batch points
			bp, _ = getPointsBatch()

			// flush cloned points
			conn.Write(sbp)
		}
	}
}

// Aggregates runtime statistics (every *aggregateInterval* period)
func dataGatherer(pointsChan chan *client.Point) {
	tags := getEnvTags(confRuntimeTagsKey)

	for range time.Tick(confMonitoringAggregationInterval) {
		/* gather mem and goroutines stats */
		memstats, err := getMemStats(tags)
		if err == nil && memstats != nil {
			pointsChan <- memstats
		}

		/* gather garbage collector stats */
		gcstats, err := getGCStats(tags)
		if err == nil && gcstats != nil {
			pointsChan <- gcstats
		}
	}
}

func getMemStats(tags map[string]string) (*client.Point, error) {
	var memStats runtime.MemStats
	runtime.ReadMemStats(&memStats)

	fields := map[string]interface{}{
		"goroutines_count": runtime.NumGoroutine(),
		"memory_alloc":     int64(memStats.Alloc),
		"cpu_count":        runtime.NumCPU(),
	}

	return client.NewPoint(runtimeTablename, tags, fields, time.Now())
}

func getGCStats(tags map[string]string) (*client.Point, error) {
	var gcStats debug.GCStats
	debug.ReadGCStats(&gcStats)

	if gcStats.LastGC.After(latestGCStart) {
		latestGCStart = gcStats.LastGC

		fields := map[string]interface{}{
			"duration_ns": gcStats.Pause[0].Nanoseconds(),
		}

		return client.NewPoint(gcTablename, tags, fields, gcStats.LastGC)
	}

	return nil, nil
}
