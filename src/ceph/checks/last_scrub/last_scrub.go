package last_scrub

import (
  "ceph"
  "time"
  "log"
  "fmt"
)

type LastScrub struct {
  lastScrub int
}

const failMessage = "PG was scrubbed recently"


func New(ls int) *LastScrub {
  check := &LastScrub{lastScrub: ls}

  return check
}


func (ls *LastScrub) Check(c *ceph.Ceph, pg ceph.PG_info) string {

   pg_stats := ceph.PG_query{}
   err := ceph.RunCephCommand(fmt.Sprintf(c.PG_query_command, pg.PG_id), &pg_stats)
   if err != nil { log.Fatal(err) }

   if time.Now().Sub(pg_stats.Info.Stats.Last_deep_scrub_stamp.Time) < time.Duration(ls.lastScrub) * time.Hour {
      return "CHECK_SKIP"
   }

   return "CHECK_OK"
}


func (ls *LastScrub) GetFailureMessage() string {
  return failMessage
}
