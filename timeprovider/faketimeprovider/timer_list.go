package faketimeprovider

import (
	"time"
)

type timeWatcher interface {
	timeUpdated(time.Time)
	completion() time.Time
}

// A PriorityQueue implements heap.Interface and holds Items.
type timerList []timeWatcher

func (tl timerList) Len() int { return len(tl) }

func (tl timerList) Less(i, j int) bool {
	return tl[i].completion().After(tl[j].completion())
}

func (tl timerList) Swap(i, j int) {
	tl[i], tl[j] = tl[j], tl[i]
}

func (tl *timerList) Push(x interface{}) {
	item := x.(timeWatcher)
	*tl = append(*tl, item)
}

func (tl *timerList) Pop() interface{} {
	old := *tl
	n := len(old)
	item := old[n-1]
	*tl = old[0 : n-1]
	return item
}
