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
}

// embed this into Ceph struct above?
type Settings struct {
  Ceph_binary string
  PG_list_stale uint

  PG_list_command string
}

type CephCheck interface {
  Check(*Ceph, PG_info) int
}

// TODO:
// proper checker functions
// use params from config file
// use proper logging
//

func New(settings Settings) *Ceph {
  ceph := &Ceph{Binary_path: settings.Ceph_binary}

  ceph.PG_list_command = ceph.Binary_path + " pg ls --format json"
  ceph.PG_state_command = ceph.Binary_path + " pg stat --format json"
  ceph.PG_query_command = ceph.Binary_path + " pg %s query"
  ceph.Health_detail_command = ceph.Binary_path + " health detail --format json"
  ceph.Deep_scrub_command = ceph.Binary_path + " pg deep-scrub %s"

  ceph.PG_list_stale = time.Duration(settings.PG_list_stale) * time.Second

  return ceph
}
