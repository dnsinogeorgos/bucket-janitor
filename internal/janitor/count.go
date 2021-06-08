package janitor

func (j *Janitor) countTypes() {
	defer j.wg.Done()

	for ts := range j.typeSetChan {
		if v, ok := j.counter.Load(ts); ok {
			i := v.(*uint64)
			*i++
		} else {
			i := new(uint64)
			*i++
			j.counter.Store(ts, i)
		}
	}
}
