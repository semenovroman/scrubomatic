package main

import (
	"ceph"
	"github.com/spf13/viper"
)

// TODO:
// have a 'current speed' that is constantly updated by goroutine ?
// deep-scrub multiple pgs on different osds\hosts?

func main() {

	viper.AddConfigPath(".")
	viper.SetConfigName("scrubomatic")
	viper.ReadInConfig()

	// viper.GetString("db_user")

	// pass all settings
	settings := ceph.Settings{
								Ceph_binary: "sudo docker exec 291 /usr/bin/ceph",
								PG_list_stale: 15,
							}

	scrubomatic := ceph.New(settings)
	scrubomatic.DeepScrub()
}
