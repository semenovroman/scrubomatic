package pgs

import (
    "ceph"
    "ceph/checks/health"
    "ceph/checks/last_scrub"
    "ceph/checks/last_change"
    "ceph/checks/io"
    "ceph/checks/concurrent_scrubs"

    "time"
    "log"
    "fmt"
    "sort"
)


func DeepScrub(c *ceph.Ceph) {

    // TEMP HACK
    c.Checks_map = make(map[string]ceph.CephCheck)
    c.Checks_map["health"] = health.New(c.Health_status)
    c.Checks_map["last_srub"] = last_scrub.New(c.Last_scrub)
    c.Checks_map["last_change"] = last_change.New(c.Last_change)
    c.Checks_map["io"] = io.New(uint64(c.Io_writes), uint64(c.Io_reads), uint64(c.Io_ops))
    c.Checks_map["concurrent_scrubs"] = concurrent_scrubs.New(uint(c.Concurrent_scrubs))

    for {
        pgs_list := GetPGsList(c)

        for _, pg := range pgs_list {
            if time.Now().Sub(c.Last_pg_list_update) > c.PG_list_stale {
                fmt.Printf("PG list is stale, refreshing... %v\n", c.Last_pg_list_update)

                break
            }

            if Check_pg(c, pg) {
                Scrub_pg(c, pg)
            }
        }
    }
}


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

    fmt.Printf("Deep-scrubbing: %s last deep-scrub on [%s] (%v ago)\n", string(pg.PG_id), last_deep_scrub, time.Now().Sub(pg.Last_deep_scrub_stamp.Time))

    err := ceph.RunCephCommand(fmt.Sprintf(c.Deep_scrub_command, pg.PG_id), nil)
    if err != nil { log.Fatal(err) }


    pgInfo := ceph.PG_query{};

    for {
        time.Sleep(500 * time.Millisecond)

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

    for name, check := range c.Checks_map {

        switch cr, err := check.Check(c, pg); cr {
            case "CHECK_WAIT": {
                time.Sleep(5 * time.Second)

                fmt.Printf(err.Error())

                for ; cr != "CHECK_OK"; {

                    switch cr, err = check.Check(c, pg); cr {
                        case "CHECK_WAIT": {
			                fmt.Printf("Check %s failed on pg %s [%s]\n", name, pg.PG_id, cr)
                            time.Sleep(5 * time.Second)
                        }
                        case "CHECK_SKIP": {
			                fmt.Printf("Check \"%s\" failed on pg %s - [%s]\n", name, pg.PG_id, cr)
                            return false
                        }
			            case "CHECK_OK": {
			            }
			            default: {
			                     log.Fatal("!!! UNKNOWN status (%s) returned by check \"%s\"", cr, name)
			            }
                    }
                }
            }
            case "CHECK_SKIP": {
		          fmt.Printf("Check %s failed on pg %s [%s]\n", name, pg.PG_id, cr)
                  return false
            }
	        case "CHECK_OK": {
	        }
	        default: {
		          log.Fatal("!!! UNKNOWN status (%s) returned by check \"%s\"", cr, name)
	        }
        }

    }

    return true
}
