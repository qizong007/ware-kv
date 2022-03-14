package util

import "math/rand"

const (
	maxLevel    = 16   // Should be enough for 2^16 elements
	probability = 0.25 // It's the best practice value (Based on Time and Space)
	// Higher level, smaller probability
	// Just modify the probability factor P to reduce additional space
)

type slNode struct {
	score float64
	val   interface{}
	// forward[i]: the slNode is current node's next one, which is at i_th layer
	forward []*slNode
}

type SlElement struct {
	Score float64     `json:"score"`
	Val   interface{} `json:"val"`
}

func newSLNode(score float64, value interface{}, level int) *slNode {
	return &slNode{
		score:   score,
		val:     value,
		forward: make([]*slNode, level),
	}
}

type SkipList struct {
	header *slNode
	len    int
	level  int
}

func NewSkipList() *SkipList {
	return &SkipList{
		header: &slNode{forward: make([]*slNode, maxLevel)},
	}
}

func randomLevel() int {
	level := 1
	for rand.Float32() < probability && level < maxLevel {
		level++
	}
	return level
}

func (sl *SkipList) Front() *slNode {
	return sl.header.forward[0]
}

func (n *slNode) Next() *slNode {
	if n != nil {
		return n.forward[0]
	}
	return nil
}

func (sl *SkipList) GetList() []SlElement {
	list := make([]SlElement, sl.len)
	x := sl.Front()
	i := 0
	for x != nil {
		list[i] = SlElement{
			Score: x.score,
			Val:   x.val,
		}
		x = x.Next()
		i++
	}
	return list
}

func (sl *SkipList) Search(score float64) (*slNode, bool) {
	x := sl.header
	for i := sl.level - 1; i >= 0; i-- {
		for x.forward[i] != nil && x.forward[i].score < score {
			x = x.forward[i]
		}
	}
	x = x.Next()
	if x != nil && x.score == score {
		return x, true
	}
	return nil, false
}

func (sl *SkipList) Insert(score float64, value interface{}) *slNode {
	update := make([]*slNode, maxLevel)
	x := sl.header
	for i := sl.level - 1; i >= 0; i-- {
		for x.forward[i] != nil && x.forward[i].score < score {
			x = x.forward[i]
		}
		update[i] = x
	}
	x = x.Next()

	if x != nil && x.score == score {
		x.val = value
		return x
	}

	level := randomLevel()
	if level > sl.level {
		level = sl.level + 1
		update[sl.level] = sl.header
		sl.level = level
	}
	newNode := newSLNode(score, value, level)
	for i := 0; i < level; i++ {
		newNode.forward[i] = update[i].forward[i]
		update[i].forward[i] = newNode
	}
	sl.len++
	return newNode
}

func (sl *SkipList) Delete(score float64) *slNode {
	update := make([]*slNode, maxLevel)
	x := sl.header
	for i := sl.level - 1; i >= 0; i-- {
		for x.forward[i] != nil && x.forward[i].score < score {
			x = x.forward[i]
		}
		update[i] = x
	}
	x = x.Next()

	if x != nil && x.score == score {
		for i := 0; i < sl.level; i++ {
			if update[i].forward[i] != x {
				return nil
			}
			update[i].forward[i] = x.forward[i]
		}
		sl.len--
		return x
	}

	return nil
}

func (sl *SkipList) Len() int {
	return sl.len
}
