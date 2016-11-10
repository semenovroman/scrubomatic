package health

import (
  "ceph"
  "log"
  "strings"
  "errors"
)

type Health struct {
  expected_status string
}

func New(status string) *Health {
  health := &Health{expected_status: status}

  return health
}

func (h *Health) Check(c *ceph.Ceph, pg ceph.PG_info) (string, error) {
  health := ceph.Ceph_health{}
  err := ceph.RunCephCommand(c.Health_detail_command, &health)
  if err != nil { log.Fatal(err) }

  if strings.Compare(health.Overall_status, h.expected_status) != 0 {
    return "CHECK_WAIT", errors.New("Cluster health status")
  }

  return "CHECK_OK", nil
}
