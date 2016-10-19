package io

import (
  "ceph"
  "log"
)

type IO struct {
  writes uint64
  reads uint64
  iops uint64
}

const failMessage = "Cluster is under heavy io"

func New(writes uint64, reads uint64, iops uint64) *IO {
  check := &IO{
                writes: writes,
                reads: reads,
                iops: iops,
              }

  return check
}

func (io *IO) Check(c *ceph.Ceph, pg ceph.PG_info) string {
  pgs_scrubbing := ceph.PGS_state{}
  err := ceph.RunCephCommand(c.PG_state_command, &pgs_scrubbing)
  if err != nil { log.Fatal(err) }

  if pgs_scrubbing.Write_bytes_sec > io.writes {
    return "CHECK_WAIT"
  }

  if pgs_scrubbing.Read_bytes_sec > io.reads {
    return "CHECK_WAIT"
  }

  if pgs_scrubbing.Io_sec > io.iops {
    return "CHECK_WAIT"
  }

  return "CHECK_OK"
}

func (io *IO) GetFailureMessage() string {
  return failMessage
}
