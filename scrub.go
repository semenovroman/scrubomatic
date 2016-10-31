package main

import (
	"ceph"
	"ceph/pgs"
	"github.com/spf13/viper"
)

// TODO:
// have a 'current speed' that is constantly updated by goroutine ?
// deep-scrub multiple pgs on different osds\hosts?

func main() {

	viper.AddConfigPath(".")
	viper.SetConfigName("scrubomatic")
	viper.ReadInConfig()

	settings := ceph.Settings{
								Ceph_binary: viper.GetString("ceph_binary"),
								PG_list_stale: viper.GetInt("pg_list_stale"),
								Health_status: viper.GetString("checks.health.status"),
								Last_scrub: viper.GetInt("checks.last_scrub.hours"),
								Last_change: viper.GetInt("checks.last_change.minutes"),
								Io_reads: viper.GetInt("checks.io.reads"),
								Io_writes: viper.GetInt("checks.io.writes"),
								Io_ops: viper.GetInt("checks.io.ops"),
								Concurrent_scrubs: viper.GetInt("concurrent_scrubs"),
							}

	scrubomatic := ceph.New(settings)
	pgs.DeepScrub(scrubomatic)
}
