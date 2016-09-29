package ceph


type ceph_health struct {
	Overall_status string
}

type pgs_state struct {
	Write_bytes_sec uint64
	Read_bytes_sec uint64
	Io_sec uint64
	Num_pg_by_state []struct {
		Name string
		Num int
	}
}

type pg_query struct {
	State string
	Info struct {
		PG_id string `json:"pgid"`
		Stats struct {
			Last_change CephTime
			Last_deep_scrub_stamp CephTime
		}
	}
}

type PG_info struct {
	PG_id string `json:"pgid"`
	Last_change CephTime
	Last_deep_scrub_stamp CephTime
	Acting_primary int
}


type PGSByDate []PG_info

func (pgs PGSByDate) Len() int {
	 return len(pgs)
}

func (pgs PGSByDate) Swap(i, j int) {
	pgs[i], pgs[j] = pgs[j], pgs[i]
}

func (pgs PGSByDate) Less(i, j int) bool {
  t1 := pgs[i].Last_deep_scrub_stamp.Time
  t2 := pgs[j].Last_deep_scrub_stamp.Time

	return t1.Before(t2)
}
