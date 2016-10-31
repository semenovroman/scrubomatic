package last_change

import (
  "ceph"
  "time"
  "log"
  "fmt"
)

type LastChange struct {
  lastChange int
}

const failMessage = "PG was changed recently"

func New(lc int) *LastChange {
  check := &LastChange{lastChange: lc}

  return check
}

func (lc *LastChange) Check(c *ceph.Ceph, pg ceph.PG_info) string {

   pg_stats := ceph.PG_query{}
   err := ceph.RunCephCommand(fmt.Sprintf(c.PG_query_command, pg.PG_id), &pg_stats)
   if err != nil { log.Fatal(err) }

   if time.Now().Sub(pg_stats.Info.Stats.Last_change.Time) < time.Duration(lc.lastChange) * time.Minute {
      return "CHECK_SKIP"
   }

   return "CHECK_OK"
}

func (lc *LastChange) GetFailureMessage() string {
  return failMessage
}
