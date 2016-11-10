package last_scrub

import (
  "ceph"
  "time"
  "log"
  "fmt"
  "errors"
)

type LastScrub struct {
  lastScrub int
}

func New(ls int) *LastScrub {
  check := &LastScrub{lastScrub: ls}

  return check
}


func (ls *LastScrub) Check(c *ceph.Ceph, pg ceph.PG_info) (string, error) {

   pg_stats := ceph.PG_query{}
   err := ceph.RunCephCommand(fmt.Sprintf(c.PG_query_command, pg.PG_id), &pg_stats)
   if err != nil { log.Fatal(err) }

   if time.Now().Sub(pg_stats.Info.Stats.Last_deep_scrub_stamp.Time) < time.Duration(ls.lastScrub) * time.Hour {
      return "CHECK_SKIP", errors.New("Last scrub")
   }

   return "CHECK_OK", nil
}
