package health

import (
  "ceph"
  "log"
)

type Health struct {
  good_status string
}

const failMessage = "Cluster is not healthy, not scrubbing"

func New(status string) *Health {
  health := &Health{good_status: status}

  return health
}

func (h *Health) Check(c *ceph.Ceph, pg ceph.PG_info) string {
  health := ceph.Ceph_health{}
  err := ceph.RunCephCommand(c.Health_detail_command, &health)
  if err != nil { log.Fatal(err) }

  if health.Overall_status != h.good_status {
    return "CHECK_WAIT"
  }

  return "CHECK_OK"
}

func (h *Health) GetFailureMessage() string {
  return failMessage
}
