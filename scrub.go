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
	scrubomatic := ceph.New("sudo docker exec c3d /usr/bin/ceph")

	scrubomatic.DeepScrub()
}
