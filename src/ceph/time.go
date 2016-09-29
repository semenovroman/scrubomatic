package ceph

import (
  "time"
  "strings"
)

const cephTimeFormat = "2006-01-02 15:04:05.000000"

type CephTime struct {
    time.Time
}

func (ct *CephTime) UnmarshalJSON(timeString []byte) (err error) {
    loc, _ := time.LoadLocation("Local")

    // can do this with bytes ( b = b[1 : len(b)-1] )
    ceph_timestamp := strings.Trim(string(timeString), "\"")

    ct.Time, err = time.ParseInLocation(cephTimeFormat, ceph_timestamp, loc)

    return err
}
