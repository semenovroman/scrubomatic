package ceph

import (
  "time"
  "fmt"
  "log"
  "sort"
)

type Ceph struct {
  binary_path string
  last_pg_list_update time.Time
  health_detail_command string
  pg_list_command string
  pgs_state_command string
  pg_query_command string
  deep_scrub_command string
  pg_list_stale time.Duration
}

// embed this into Ceph struct above?
type Settings struct {
  Ceph_binary string
  PG_list_stale uint
}

// TODO:
// proper checker functions
// use params from config file
// use proper logging
//

func New(settings Settings) *Ceph {
  ceph := &Ceph{binary_path: settings.Ceph_binary}

  ceph.pg_list_command = ceph.binary_path + " pg ls --format json"
  ceph.pgs_state_command = ceph.binary_path + " pg stat --format json"
  ceph.pg_query_command = ceph.binary_path + " pg %s query"
  ceph.health_detail_command = ceph.binary_path + " health detail --format json"
  ceph.deep_scrub_command = ceph.binary_path + " pg deep-scrub %s"

  ceph.pg_list_stale = time.Duration(settings.PG_list_stale) * time.Second

  return ceph
}

func (ceph *Ceph) DeepScrub() {
  for {
    pgs_list := ceph.GetPGsList()

    for _, pg := range pgs_list {

      if time.Now().Sub(ceph.last_pg_list_update) > ceph.pg_list_stale {
        fmt.Printf("PG list is stale, refreshing... %v\n", ceph.last_pg_list_update)
        time.Sleep(5 * time.Second)
        
        break
      }

      if ceph.Check_pg(pg) {
        ceph.Scrub_pg(pg)
      }
    }
  }
}


func (ceph *Ceph) GetHealth() string {

  health := ceph_health{}
  err := runCephCommand(ceph.health_detail_command, &health)
  if err != nil { log.Fatal(err) }

  return health.Overall_status
}


func (ceph *Ceph) GetPGsList() PGSByDate {

  pg_list := PGSByDate{}
  ok := runCephCommand(ceph.pg_list_command, &pg_list)
  if ok != nil { log.Fatal(ok) }

  sort.Sort(pg_list)

  ceph.last_pg_list_update = time.Now()

  return pg_list
}


func (ceph *Ceph) Check_pg(pg PG_info) bool {

	if time.Now().Sub(pg.Last_deep_scrub_stamp.Time) / time.Hour < -999 {
		fmt.Printf("Scrubbed less than 24 hrs before, not scrubbing\n")

		return false
	}

	if time.Now().Sub(pg.Last_change.Time) / time.Minute < -100000 {
		fmt.Printf("Placement group was recently written to, not scrubbing\n")

		return false
	}

	if ceph.GetHealth() != "HEALTH_WARN" {
		fmt.Printf("Cluster is not healthy, not scrubbing\n")

		return false
	}

  // get number of pgs in different state (and speed) to determine if we have any currently scrubbing
	pgs_scrubbing := pgs_state{}
  ok := runCephCommand(ceph.pgs_state_command, &pgs_scrubbing)
  if ok != nil { log.Fatal(ok) }

  for _, pgs := range pgs_scrubbing.Num_pg_by_state {
    if pgs.Name == "active+clean+deep-scrubbing" {
      if pgs.Num > 0 {
        return false
      }
    }
  }

	if pgs_scrubbing.Write_bytes_sec > (1024 * 1024 * 1) {
		fmt.Printf("Current write traffic > 1MB\n")

		return false
	}

	if pgs_scrubbing.Read_bytes_sec > (1024 * 1024 * 1) {
		fmt.Printf("Current read traffic > 1MB\n")

		return false
	}

	if pgs_scrubbing.Io_sec > (1024 * 1024) {
		fmt.Printf("Current iops > 1 million\n")

		return false
	}


	return true
}


func (ceph *Ceph) Scrub_pg(pg PG_info) {

  last_deep_scrub := pg.Last_deep_scrub_stamp
  deep_scrub_start := time.Now()

  fmt.Printf("Deep-scrubbing: %s %s\n", string(pg.PG_id), last_deep_scrub)

  err := runCephCommand(fmt.Sprintf(ceph.deep_scrub_command, pg.PG_id), nil)
  if err != nil { log.Fatal(err) }


  pgInfo := pg_query{};

  for {
    time.Sleep(100 * time.Millisecond)

    err = runCephCommand(fmt.Sprintf(ceph.pg_query_command, pg.PG_id), &pgInfo)
    if err != nil { log.Fatal(err) }

    // diff := deep_scrub_start.Sub(pgInfo.Info.Stats.Last_deep_scrub_stamp.Time)

    if !last_deep_scrub.Equal(pgInfo.Info.Stats.Last_deep_scrub_stamp.Time) {
      fmt.Printf("Finished deep-scrubbing %s in %v\n", pg.PG_id, time.Now().Sub(deep_scrub_start))
      break
    }

    if time.Now().Sub(deep_scrub_start) > 180 * time.Second {
      fmt.Printf("Timed out waiting for deep-scrub to complete\n")
      break
    }

    time.Sleep(1 * time.Second)
  }
}
