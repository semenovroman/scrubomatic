package health

import (
  "ceph"
  "log"
  "strings"
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

  if strings.Compare(health.Overall_status, h.good_status) != 0 {
    return "CHECK_WAIT"
  }

  return "CHECK_OK"
}

func (h *Health) GetFailureMessage() string {
  return "Cluster health status is not OK: " + h.good_status
}
