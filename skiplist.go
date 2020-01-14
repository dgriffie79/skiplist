package skiplist

import (
	"math"
	"math/rand"
	"unsafe"
)

const (
	MAX_HEIGHT = 63
)

type node struct {
	key   uint64
	right *node
	down  *node
}

type leaf struct {
	key   uint64
	right *node
	val   interface{}
}

type finger struct {
	list *Skiplist
	pred [MAX_HEIGHT]*node
}

type Skiplist struct {
	head  [MAX_HEIGHT]*node
	sent  *node
	level int
	P     float64
}

func NewSkiplist() *Skiplist {
	var l Skiplist

	l.sent = &node{math.MaxUint64, nil, nil}
	l.sent.right = l.sent

	l.head[0] = &node{0, l.sent, nil}
	for i := 1; i < MAX_HEIGHT; i++ {
		l.head[i] = &node{0, l.sent, l.head[i-1]}
	}

	l.level = 0
	l.P = 1 / math.E

	return &l
}

func (l *Skiplist) Get(key uint64) interface{} {
	n := l.head[l.level]
	for level := l.level; level > 0; level-- {
		for key > n.right.key {
			n = n.right
		}
		n = n.down
	}
	for key > n.right.key {
		n = n.right
	}

	if key == n.right.key {
		return (*leaf)(unsafe.Pointer(n.right)).val
	}
	return nil
}

func (l *Skiplist) Set(key uint64, val interface{}) {
	var pred [MAX_HEIGHT]*node

	for n, level := l.head[l.level], l.level; level >= 0; level-- {
		for key > n.right.key {
			n = n.right
		}
		pred[level] = n
		n = n.down
	}

	if pred[0].right.key == key {
		(*leaf)(unsafe.Pointer(pred[0].right)).val = val
		return
	}

	level := 0
	for rand.Float64() < l.P {
		level++
	}

	for l.level < level {
		pred[l.level+1] = l.head[l.level+1]
		l.level++
	}

	n := (*node)(unsafe.Pointer(&leaf{key, pred[0].right, val}))
	pred[0].right = n
	for i := 1; i <= level; i++ {
		nn := &node{key, pred[i].right, n}
		pred[i].right = nn
		n = nn
	}
}

func (l *Skiplist) Del(key uint64) {
	n := l.head[l.level]
	for level := l.level; level > 0; level-- {
		for key > n.right.key {
			n = n.right
		}
		if key == n.right.key {
			n.right = n.right.right
		}
		n = n.down
	}
	for key > n.right.key {
		n = n.right
	}
	if key == n.right.key {
		n.right = n.right.right
	}
}

func (l *Skiplist) Finger() *finger {
	return &finger{l, l.head}
}

func (f *finger) Reset() {
	f.pred = f.list.head
}

func (f *finger) Get(key uint64) interface{} {
	level := f.list.level
	for level > 0 {
		for key > f.pred[level].right.key {
			f.pred[level] = f.pred[level].right
		}
		f.pred[level-1] = f.pred[level].down
		level--
	}
	for key > f.pred[0].right.key {
		f.pred[0] = f.pred[0].right
	}

	if key == f.pred[0].right.key {
		return (*leaf)(unsafe.Pointer(f.pred[0].right)).val
	}
	return nil
}

func (f *finger) Next() (uint64, interface{}) {
	f.pred[0] = f.pred[0].right
	return f.pred[0].key, (*leaf)(unsafe.Pointer(f.pred[0])).val
}
