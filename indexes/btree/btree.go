package btree

import (
	"github.com/shamcode/simd/storage"
)

// Based on
// https://github.com/emirpasic/gods/tree/master/trees/btree
// https://github.com/google/btree/blob/master/btree.go

type Key interface {
	Less(than Key) bool
}

type entry struct {
	key     Key
	records storage.IDStorage
}

type node struct {
	parent   *node
	entries  []*entry
	children []*node
}

func (node *node) search(key Key) (index int, found bool) {
	low, high := 0, len(node.entries)-1
	var mid int
	for low <= high {
		mid = int(uint(high+low) >> 1)
		itemKey := node.entries[mid].key
		if key.Less(itemKey) {
			high = mid - 1
		} else if itemKey.Less(key) {
			low = mid + 1
		} else {
			return mid, true
		}
	}
	return low, false
}

func (node *node) iterateAscend(start, stop Key, includeStart bool, hit bool, iter func(e *entry)) (bool, bool) {
	var ok bool
	var index int
	if nil != start {
		index, _ = node.search(start)
	}
	for i := index; i < len(node.entries); i++ {
		if len(node.children) > 0 {
			if hit, ok = node.children[i].iterateAscend(start, stop, includeStart, hit, iter); !ok {
				return hit, false
			}
		}
		if !includeStart && !hit && start != nil && !start.Less(node.entries[i].key) {
			hit = true
			continue
		}
		hit = true
		if stop != nil && !node.entries[i].key.Less(stop) {
			return hit, false
		}
		iter(node.entries[i])
	}
	if len(node.children) > 0 {
		if hit, ok = node.children[len(node.children)-1].iterateAscend(start, stop, includeStart, hit, iter); !ok {
			return hit, false
		}
	}
	return hit, true
}

func (node *node) iterateDescend(start, stop Key, includeStart bool, hit bool, iter func(e *entry)) (bool, bool) {
	var ok, found bool
	var index int
	if start != nil {
		index, found = node.search(start)
		if !found {
			index = index - 1
		}
	} else {
		index = len(node.entries) - 1
	}
	for i := index; i >= 0; i-- {
		if start != nil && !node.entries[i].key.Less(start) {
			if !includeStart || hit || start.Less(node.entries[i].key) {
				continue
			}
		}
		if len(node.children) > 0 {
			if hit, ok = node.children[i+1].iterateDescend(start, stop, includeStart, hit, iter); !ok {
				return hit, false
			}
		}
		if stop != nil && !stop.Less(node.entries[i].key) {
			return hit, false
		}
		hit = true
		iter(node.entries[i])
	}
	if len(node.children) > 0 {
		if hit, ok = node.children[0].iterateDescend(start, stop, includeStart, hit, iter); !ok {
			return hit, false
		}
	}
	return hit, true
}

var _ BTree = (*btree)(nil)

type btree struct {
	root        *node
	maxChildren int
}

func (tree *btree) Get(indexKey interface{}) storage.IDStorage {
	key, ok := indexKey.(Key)
	if !ok {
		return nil
	}
	return tree.GetForKey(key)
}

func (tree *btree) GetForKey(key Key) storage.IDStorage {
	node, index, found := tree.searchRecursively(tree.root, key)
	if found {
		return node.entries[index].records
	}
	return nil
}

func (tree *btree) Set(indexKey interface{}, records storage.IDStorage) {
	key, ok := indexKey.(Key)
	if ok {
		tree.SetForKey(key, records)
	}
}

func (tree *btree) SetForKey(key Key, records storage.IDStorage) {
	e := &entry{key: key, records: records}
	if tree.root == nil {
		tree.root = &node{
			entries:  []*entry{e},
			children: []*node{},
		}
	} else {
		tree.insert(tree.root, e)
	}
}

func (tree *btree) isLeaf(node *node) bool {
	return len(node.children) == 0
}

func (tree *btree) middle() int {
	return (tree.maxChildren - 1) / 2
}

func (tree *btree) maxEntries() int {
	return tree.maxChildren - 1
}

func (tree *btree) searchRecursively(startNode *node, key Key) (node *node, index int, found bool) {
	if nil == tree.root {
		return nil, -1, false
	}
	node = startNode
	for {
		index, found = node.search(key)
		if found {
			return node, index, true
		}
		if tree.isLeaf(node) {
			return nil, -1, false
		}
		node = node.children[index]
	}
}

func (tree *btree) insert(node *node, entry *entry) {
	if tree.isLeaf(node) {
		tree.insertIntoLeaf(node, entry)
	} else {
		tree.insertIntoInternal(node, entry)
	}
}

func (tree *btree) insertIntoLeaf(node *node, entry *entry) {
	insertPosition, found := node.search(entry.key)
	if !found {
		// Insert entry's key in the middle of the node
		node.entries = append(node.entries, nil)
		copy(node.entries[insertPosition+1:], node.entries[insertPosition:])
		node.entries[insertPosition] = entry
		tree.split(node)
	}
}

func (tree *btree) insertIntoInternal(node *node, entry *entry) {
	insertPosition, found := node.search(entry.key)
	if !found {
		tree.insert(node.children[insertPosition], entry)
	}
}

func (tree *btree) split(node *node) {
	if len(node.entries) <= tree.maxEntries() {
		return
	}
	if node == tree.root {
		tree.splitRoot()
	} else {
		tree.splitNonRoot(node)
	}
}

func (tree *btree) splitNonRoot(n *node) {
	middle := tree.middle()
	parent := n.parent

	left := &node{entries: append([]*entry(nil), n.entries[:middle]...), parent: parent}
	right := &node{entries: append([]*entry(nil), n.entries[middle+1:]...), parent: parent}

	// Move children from the node to be split into left and right nodes
	if !tree.isLeaf(n) {
		left.children = append([]*node(nil), n.children[:middle+1]...)
		right.children = append([]*node(nil), n.children[middle+1:]...)
		setParent(left.children, left)
		setParent(right.children, right)
	}

	insertPosition, _ := parent.search(n.entries[middle].key)

	// Insert middle key into parent
	parent.entries = append(parent.entries, nil)
	copy(parent.entries[insertPosition+1:], parent.entries[insertPosition:])
	parent.entries[insertPosition] = n.entries[middle]

	// Set child left of inserted key in parent to the created left node
	parent.children[insertPosition] = left

	// Set child right of inserted key in parent to the created right node
	parent.children = append(parent.children, nil)
	copy(parent.children[insertPosition+2:], parent.children[insertPosition+1:])
	parent.children[insertPosition+1] = right

	tree.split(parent)
}

func (tree *btree) splitRoot() {
	middle := tree.middle()

	left := &node{entries: append([]*entry(nil), tree.root.entries[:middle]...)}
	right := &node{entries: append([]*entry(nil), tree.root.entries[middle+1:]...)}

	if !tree.isLeaf(tree.root) {
		left.children = append([]*node(nil), tree.root.children[:middle+1]...)
		right.children = append([]*node(nil), tree.root.children[middle+1:]...)
		setParent(left.children, left)
		setParent(right.children, right)
	}
	newRoot := &node{
		entries:  []*entry{tree.root.entries[middle]},
		children: []*node{left, right},
	}

	left.parent = newRoot
	right.parent = newRoot
	tree.root = newRoot
}

type direction uint8

const (
	descend direction = iota + 1
	ascend
)

func (tree *btree) collect(dir direction, start, stop Key, includeStart bool, hit bool) (int, []storage.IDIterator) {
	if nil == tree.root {
		return 0, nil
	}
	var count int
	var ids []storage.IDIterator
	iter := func(e *entry) {
		itemCount := e.records.Count()
		if itemCount > 0 {
			count += itemCount
			ids = append(ids, e.records)
		}
	}
	switch dir {
	case ascend:
		tree.root.iterateAscend(start, stop, includeStart, hit, iter)
	case descend:
		tree.root.iterateDescend(start, stop, includeStart, hit, iter)
	}
	return count, ids
}

func (tree *btree) LessThan(key Key) (int, []storage.IDIterator) {
	return tree.collect(ascend, nil, key, false, false)
}

func (tree *btree) LessOrEqual(key Key) (int, []storage.IDIterator) {
	return tree.collect(descend, key, nil, true, false)
}

func (tree *btree) GreaterThan(key Key) (int, []storage.IDIterator) {
	return tree.collect(descend, nil, key, false, false)
}

func (tree *btree) GreaterOrEqual(key Key) (int, []storage.IDIterator) {
	return tree.collect(ascend, key, nil, true, false)
}

func (tree *btree) All(iter func(key Key, records storage.IDStorage)) {
	if nil == tree.root {
		return
	}
	tree.root.iterateAscend(nil, nil, false, false, func(e *entry) {
		iter(e.key, e.records)
	})
}

func (tree *btree) ForKey(key Key) (int, storage.IDIterator) {
	idStorage := tree.GetForKey(key)
	if nil == idStorage {
		return 0, nil
	}
	return idStorage.Count(), idStorage
}

func setParent(nodes []*node, parent *node) {
	for _, node := range nodes {
		node.parent = parent
	}
}

func NewTree(maxChildren int, uniq bool) BTree {
	return &btree{
		maxChildren: maxChildren,
	}
}
