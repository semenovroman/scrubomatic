package last_change

import (
  "ceph"
  "time"
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
   if time.Now().Sub(pg.Last_change.Time) / time.Minute < time.Duration(lc.lastChange) * time.Minute {
      return "CHECK_SKIP"
   }

   return "CHECK_OK"
}

func (lc *LastChange) GetFailureMessage() string {
  return failMessage
}
