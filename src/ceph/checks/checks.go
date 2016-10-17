package checks

import (
  "log"
  "time"
  "ceph"
)

func HealthCheck(c *ceph.Ceph) string {

  health := ceph.Ceph_health{}
  err := ceph.RunCephCommand(c.Health_detail_command, &health)
  if err != nil { log.Fatal(err) }

  return health.Overall_status
}

func LastDeepScrubCheck(c *ceph.Ceph, pg ceph.PG_info) time.Duration {
  return time.Now().Sub(pg.Last_deep_scrub_stamp.Time) / time.Hour
}

func LastPGChangeCheck(c *ceph.Ceph, pg ceph.PG_info) time.Duration {
  return time.Now().Sub(pg.Last_change.Time) / time.Minute
}

func WritesCheck(c *ceph.Ceph, pgs_scrubbing ceph.PGS_state) uint64 {
  return pgs_scrubbing.Write_bytes_sec
}

func ReadsCheck(c *ceph.Ceph, pgs_scrubbing ceph.PGS_state) uint64 {
  return pgs_scrubbing.Read_bytes_sec
}

func IopsCheck(c *ceph.Ceph, pgs_scrubbing ceph.PGS_state) uint64 {
  return pgs_scrubbing.Io_sec
}
