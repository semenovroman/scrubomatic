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

var checks_map map[string]ceph.CephCheck

func DeepScrub(c *ceph.Ceph) {

    checks_map := make(map[string]ceph.CephCheck)
    checks_map["health"] = health.New("HEALTH_WARN")
    checks_map["last_srub"] = last_scrub.New(-999)
    checks_map["last_change"] = last_change.New(-999)
    checks_map["io"] = io.New(1024 * 1024, 1024 * 1024, 1024 * 1024)
    checks_map["concurrent_scrubs"] = concurrent_scrubs.New(0)

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
    
    for name, check := range checks_map {
        switch cr := check.Check(c, pg); cr {
            case "CHECK_WAIT": {
                fmt.Println(name + " " + check.GetFailureMessage())
                time.Sleep(3 * time.Second)

                for {
                    switch rerun := check.Check(c, pg); rerun {
                        case "CHECK_WAIT": {
                            time.Sleep(3 * time.Second)
                        }
                        case "CHECK_SKIP": {
                            return false
                        }
                    }
                }
            }
            case "CHECK_SKIP": {
                fmt.Println(name + " skipping pg")
                return false
            }
        }
    }

    return true
}
