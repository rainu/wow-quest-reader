package processor

type idIterator struct {
	curId        int64
	maxId        int64
	alreadyKnown map[int64]bool
}

func newIter(knownIds []int64) *idIterator {
	result := &idIterator{
		curId:        0,
		maxId:        100000,
		alreadyKnown: map[int64]bool{},
	}

	for _, id := range knownIds {
		result.alreadyKnown[id] = true
	}

	return result
}

func (i *idIterator) Next() int64 {
	//find next id...
	for i.alreadyKnown[i.curId] {
		i.curId++
	}

	if i.curId > i.maxId {
		return -1 //end is reached
	}

	next := i.curId
	i.curId++

	return next
}
