package processor

type qIdIterator struct {
	curId        int64
	maxId        int64
	alreadyKnown map[int64]bool
}

func newQuestIter(knownIds []int64) *qIdIterator {
	result := &qIdIterator{
		curId:        0,
		maxId:        80000,
		alreadyKnown: map[int64]bool{},
	}

	for _, id := range knownIds {
		result.alreadyKnown[id] = true
	}

	return result
}

func (i *qIdIterator) Next() int64 {
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
