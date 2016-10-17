package pgs

import (
  "ceph"
  "ceph/checks"

  "time"
  "log"
  "fmt"
  "sort"
)

func GetPGsList(c *ceph.Ceph) ceph.PGSByDate {

  pg_list := ceph.PGSByDate{}
  ok := ceph.RunCephCommand(c.PG_list_command, &pg_list)
  if ok != nil { log.Fatal(ok) }

  sort.Sort(pg_list)

  c.Last_pg_list_update = time.Now()

  return pg_list
}

func Scrub_pg(c *ceph.Ceph, pg ceph.PG_info) {

  last_deep_scrub := pg.Last_deep_scrub_stamp
  deep_scrub_start := time.Now()

  fmt.Printf("Deep-scrubbing: %s %s\n", string(pg.PG_id), last_deep_scrub)

  err := ceph.RunCephCommand(fmt.Sprintf(c.Deep_scrub_command, pg.PG_id), nil)
  if err != nil { log.Fatal(err) }


  pgInfo := ceph.PG_query{};

  for {
    time.Sleep(100 * time.Millisecond)

    err = ceph.RunCephCommand(fmt.Sprintf(c.PG_query_command, pg.PG_id), &pgInfo)
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

func Check_pg(c *ceph.Ceph, pg ceph.PG_info) bool {

  // for check := range checks Ceph.checks {  for check.Check() != 0; sleep check.Check()... }

	if checks.LastDeepScrubCheck(c, pg) < -999 * time.Hour {
		fmt.Printf("Scrubbed less than 24 hrs before, not scrubbing\n")

		return false
	}

	if checks.LastPGChangeCheck(c, pg) < -100000 * time.Minute {
		fmt.Printf("Placement group was recently written to, not scrubbing\n")

		return false
	}

	if checks.HealthCheck(c) != "HEALTH_WARN" {
		fmt.Printf("Cluster is not healthy, not scrubbing\n")

		return false
	}

  // get number of pgs in different state (and speed) to determine if we have any currently scrubbing
	pgs_scrubbing := ceph.PGS_state{}
  ok := ceph.RunCephCommand(c.PG_state_command, &pgs_scrubbing)
  if ok != nil { log.Fatal(ok) }

  for _, pgs := range pgs_scrubbing.Num_pg_by_state {
    if pgs.Name == "active+clean+deep-scrubbing" {
      if pgs.Num > 0 {
        return false
      }
    }
  }

	if checks.WritesCheck(c, pgs_scrubbing) > (1024 * 1024 * 1) {
		fmt.Printf("Current write traffic > 1MB\n")

		return false
	}

	if checks.ReadsCheck(c, pgs_scrubbing) > (1024 * 1024 * 1) {
		fmt.Printf("Current read traffic > 1MB\n")

		return false
	}

	if checks.IopsCheck(c, pgs_scrubbing) > (1024 * 1024) {
		fmt.Printf("Current iops > 1 million\n")

		return false
	}

	return true
}

func DeepScrub(c *ceph.Ceph) {
  for {
    pgs_list := GetPGsList(c)

    for _, pg := range pgs_list {

      if time.Now().Sub(c.Last_pg_list_update) > c.PG_list_stale {
        fmt.Printf("PG list is stale, refreshing... %v\n", c.Last_pg_list_update)
        time.Sleep(5 * time.Second)

        break
      }

      if Check_pg(c, pg) {
        Scrub_pg(c, pg)
      }
    }
  }
}
