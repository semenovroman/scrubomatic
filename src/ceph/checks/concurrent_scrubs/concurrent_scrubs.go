package concurrent_scrubs

import (
  "ceph"
  "log"
  "errors"
)

type ConcurrentScrubs struct {
  concurrentScrubs int
}

func New(cs uint) *ConcurrentScrubs {
  check := &ConcurrentScrubs{concurrentScrubs: int(cs)}

  return check
}

func (cs *ConcurrentScrubs) Check(c *ceph.Ceph, pg ceph.PG_info) (string, error) {
  pgs_scrubbing := ceph.PGS_state{}
  err := ceph.RunCephCommand(c.PG_state_command, &pgs_scrubbing)
  if err != nil { log.Fatal(err) }

  for _, pgs := range pgs_scrubbing.Num_pg_by_state {
    if pgs.Name == "active+clean+scrubbing+deep" {
      if pgs.Num > cs.concurrentScrubs {
        return "CHECK_WAIT", errors.New("More concurrent scrubs detected than allowed")
      }
    }
  }

  return "CHECK_OK", nil
}
