package last_scrub

import (
  "ceph"
  "time"
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
   if time.Now().Sub(pg.Last_deep_scrub_stamp.Time) / time.Hour < time.Duration(ls.lastScrub) * time.Hour {
      return "CHECK_SKIP"
   }

   return "CHECK_OK"
}

func (ls *LastScrub) GetFailureMessage() string {
  return failMessage
}
