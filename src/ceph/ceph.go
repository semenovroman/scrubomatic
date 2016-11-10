package ceph

import (
  "time"
)

type Ceph struct {
  Binary_path string
  Last_pg_list_update time.Time
  Health_detail_command string
  PG_list_command string
  PG_state_command string
  PG_query_command string
  Deep_scrub_command string
  PG_list_stale time.Duration

  Health_status string
  Last_scrub int
  Last_change int
  Io_reads int
  Io_writes int
  Io_ops int
  Concurrent_scrubs int

  Checks_map map[string]CephCheck
}

// embed this into Ceph struct above?
type Settings struct {
  Ceph_binary string
  PG_list_stale int

  PG_list_command string

  Health_status string
  Last_scrub int
  Last_change int
  Io_reads int
  Io_writes int
  Io_ops int
  Concurrent_scrubs int
}

type CephCheck interface {
  Check(*Ceph, PG_info) (string, error)
}

// TODO:
// use proper logging
// refactor parameter passing
// refactor check packages

func New(settings Settings) *Ceph {
  ceph := &Ceph{Binary_path: settings.Ceph_binary}

  ceph.PG_list_command = ceph.Binary_path + " pg ls --format json"
  ceph.PG_state_command = ceph.Binary_path + " pg stat --format json"
  ceph.PG_query_command = ceph.Binary_path + " pg %s query"
  ceph.Health_detail_command = ceph.Binary_path + " health detail --format json"
  ceph.Deep_scrub_command = ceph.Binary_path + " pg deep-scrub %s"

  ceph.PG_list_stale = time.Duration(settings.PG_list_stale) * time.Hour

  ceph.Health_status = settings.Health_status
  ceph.Last_scrub = settings.Last_scrub
  ceph.Last_change = settings.Last_change
  ceph.Io_reads = settings.Io_reads
  ceph.Io_writes = settings.Io_writes
  ceph.Io_ops = settings.Io_ops
  ceph.Concurrent_scrubs = settings.Concurrent_scrubs

  return ceph
}
