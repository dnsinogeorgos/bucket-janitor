package janitor

import (
	"sort"
	"sync"
)

type resultSet []result

func (rs resultSet) Len() int           { return len(rs) }
func (rs resultSet) Less(i, j int) bool { return *rs[i].Occurences < *rs[j].Occurences }
func (rs resultSet) Swap(i, j int)      { rs[i], rs[j] = rs[j], rs[i] }

type result struct {
	Mime       string
	Magic      string
	Occurences *uint64
}

func prepareResultSet(counter *sync.Map) resultSet {
	rs := resultSet{}

	counter.Range(func(key, value interface{}) bool {
		k := key.(TypeSet)
		v := value.(*uint64)

		rs = append(rs, result{k.Mime, k.Magic, v})

		return true
	})

	sort.Sort(sort.Reverse(rs))

	return rs
}
