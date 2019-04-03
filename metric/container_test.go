// Copyright 2015 bs authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package metric

import (
	"encoding/json"

	docker "github.com/fsouza/go-dockerclient"
	"gopkg.in/check.v1"
)

func (s *S) TestStatsToMetricsMap(c *check.C) {
	jsonStats := `{
       "read" : "2015-01-08T22:57:31.547920715Z",
       "network" : {
          "rx_dropped" : 0,
          "rx_bytes" : 648,
          "rx_errors" : 0,
          "tx_packets" : 8,
          "tx_dropped" : 0,
          "rx_packets" : 8,
          "tx_errors" : 0,
          "tx_bytes" : 649
       },
       "memory_stats" : {
          "stats" : {
             "total_pgmajfault" : 0,
             "cache" : 0,
             "mapped_file" : 0,
             "total_inactive_file" : 0,
             "pgpgout" : 414,
             "rss" : 6537216,
             "total_mapped_file" : 0,
             "writeback" : 0,
             "unevictable" : 0,
             "pgpgin" : 477,
             "total_unevictable" : 0,
             "pgmajfault" : 0,
             "total_rss" : 6537216,
             "total_rss_huge" : 6291456,
             "total_writeback" : 0,
             "total_inactive_anon" : 0,
             "rss_huge" : 6291456,
             "hierarchical_memory_limit": 67108864,
             "total_pgfault" : 964,
             "total_active_file" : 0,
             "active_anon" : 6537216,
             "total_active_anon" : 6537216,
             "total_pgpgout" : 414,
             "total_cache" : 0,
             "inactive_anon" : 0,
             "active_file" : 0,
             "pgfault" : 964,
             "inactive_file" : 0,
             "total_pgpgin" : 477,
             "swap" : 47312896,
             "hierarchical_memsw_limit" : 1610612736
          },
          "max_usage" : 6651904,
          "usage" : 6537216,
          "failcnt" : 0,
          "limit" : 67108864
       },
       "blkio_stats": {
          "io_service_bytes_recursive": [
             {
                "major": 8,
                "minor": 0,
                "op": "Read",
                "value": 428795731968
             },
             {
                "major": 8,
                "minor": 0,
                "op": "Write",
                "value": 388177920
             }
          ],
          "io_serviced_recursive": [
             {
                "major": 8,
                "minor": 0,
                "op": "Read",
                "value": 25994442
             },
             {
                "major": 8,
                "minor": 0,
                "op": "Write",
                "value": 1734
             }
          ],
          "io_queue_recursive": [],
          "io_service_time_recursive": [],
          "io_wait_time_recursive": [],
          "io_merged_recursive": [],
          "io_time_recursive": [],
          "sectors_recursive": []
       },
       "cpu_stats" : {
          "cpu_usage" : {
             "percpu_usage" : [
                16970827,
                1839451,
                7107380,
                10571290
             ],
             "usage_in_usermode" : 10000000,
             "total_usage" : 36488948,
             "usage_in_kernelmode" : 20000000
          },
          "system_cpu_usage" : 20091722000000000
       }
    }`
	var stats docker.Stats
	err := json.Unmarshal([]byte(jsonStats), &stats)
	c.Assert(err, check.IsNil)
	metricsMap, err := statsToMetricsMap(&stats)
	c.Assert(err, check.IsNil)
	c.Assert(metricsMap["mem_max"], check.Equals, float(6537216))
	c.Assert(metricsMap["mem_limit"], check.Equals, float(67108864))
	c.Assert(metricsMap["swap"], check.Equals, float(47312896))
	c.Assert(metricsMap["swap_limit"], check.Equals, float(1543503872))
	c.Assert(metricsMap["netrx"], check.Equals, float(648))
	c.Assert(metricsMap["nettx"], check.Equals, float(649))
	diffMemPctMax := 9.74 - metricsMap["mem_pct_max"]
	c.Assert(diffMemPctMax < 0.01, check.Equals, true)
	diffCpuMax := 0 - metricsMap["cpu_max"]
	c.Assert(diffCpuMax < 0.01, check.Equals, true)
	stats.Networks = map[string]docker.NetworkStats{
		"eth0": {
			RxBytes: 1,
			TxBytes: 2,
		},
		"eth1": {
			RxBytes: 3,
			TxBytes: 4,
		},
	}
	metricsMap, err = statsToMetricsMap(&stats)
	c.Assert(err, check.IsNil)
	c.Assert(metricsMap["netrx"], check.Equals, float(4))
	c.Assert(metricsMap["nettx"], check.Equals, float(6))
}
