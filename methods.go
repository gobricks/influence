package influence

import (
	"os"
	"strings"

	client "github.com/influxdata/influxdb/client/v2"
)

// getPointsBatch returns batch points holder
func getPointsBatch() (client.BatchPoints, error) {
	return client.NewBatchPoints(client.BatchPointsConfig{
		Database:  confDatabase,
		Precision: "ns",
	})
}

// getEnvTags parses Influx tags from environment variable with given key.
// Tags format: comma separated "key:value", e.g. "cpu:cpu-total,dc:us-west1"
func getEnvTags(key string) map[string]string {
	res := make(map[string]string)

	envVar := os.Getenv(key)
	if envVar == "" {
		return res
	}

	for _, kvPair := range strings.Split(strings.TrimSpace(envVar), ",") {
		parts := strings.Split(strings.TrimSpace(kvPair), ":")
		if len(parts) < 2 {
			continue
		}

		key, value := strings.TrimSpace(parts[0]), strings.TrimSpace(parts[1])
		if key == "" || value == "" {
			continue
		}

		res[key] = value
	}

	return res
}
