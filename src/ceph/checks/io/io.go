package io

import (
  "ceph"
  "log"
  "errors"
)

type IO struct {
  writes uint64
  reads uint64
  iops uint64
}

func New(writes uint64, reads uint64, iops uint64) *IO {

  check := &IO{ writes: writes, reads: reads, iops: iops, }

  return check
}


func (io *IO) Check(c *ceph.Ceph, pg ceph.PG_info) (string, error) {

  pgs_scrubbing := ceph.PGS_state{}
  err := ceph.RunCephCommand(c.PG_state_command, &pgs_scrubbing)
  if err != nil { log.Fatal(err) }

  if pgs_scrubbing.Write_bytes_sec > io.writes {
    return "CHECK_WAIT", errors.New("Writes")
  }

  if pgs_scrubbing.Read_bytes_sec > io.reads {
    return "CHECK_WAIT", errors.New("Reads")
  }

  if pgs_scrubbing.Io_sec > io.iops {
    return "CHECK_WAIT", errors.New("IOPS")
  }

  return "CHECK_OK", nil
}
