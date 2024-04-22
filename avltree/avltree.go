package avltree

import (
	"errors"
	"fmt"
	"strconv"
	"strings"
)

const DEPTH_INITIAL int64 = -1
const INTEGER_MAX_VALUE int64 = 9223372036854775807

type ForEachMode byte

const (
	InOrder        ForEachMode = 0
	ReverseInOrder ForEachMode = 1
	Queue          ForEachMode = 2
	Stack          ForEachMode = 3
)

type KeyCollisionBehavior byte

const (
	Replace         KeyCollisionBehavior = 0
	IgnoreInsertion KeyCollisionBehavior = 1
)

// just an alias, since the compiler can't properly parse the Comparator's generic:
// "K *any" is interpreted as a multiplication
//type pointerToAny = *any

type ToStringable interface {
	String() string
}

type KeyVal[K any, V any] struct {
	key   K
	value V
}

type KeyExtractor[K any, V any] func(value V) K

type Comparator[K any] func(k1 K, k2 K) int64

type AVLTNode[K any, V any] struct {
	keyVal       KeyVal[K, V]
	height       int64
	sizeLeft     int64
	sizeRight    int64
	father       *AVLTNode[K, V]
	left         *AVLTNode[K, V]
	right        *AVLTNode[K, V]
	nextInOrder  *AVLTNode[K, V]
	prevInOrder  *AVLTNode[K, V]
	prevInserted *AVLTNode[K, V]
	nextInserted *AVLTNode[K, V]
}

type AVLTreeConstructorParams[K any, V any] struct {
	KeyCollisionBehavior KeyCollisionBehavior
	KeyZeroValue         K
	ValueZeroValue       V
	KeyExtractor         KeyExtractor[K, V]
	Comparator           Comparator[K]
}

type AVLTree[K any, V any] struct {
	size                     int64
	avlTreeConstructorParams AVLTreeConstructorParams[K, V]
	root                     *AVLTNode[K, V]
	_NIL                     *AVLTNode[K, V]
	minValue                 *AVLTNode[K, V] // used for optimizations
	firstInserted            *AVLTNode[K, V]
}

//

// PRIVATE functions

//

func (t *AVLTree[K, V]) newNode(key K, value V) *AVLTNode[K, V] {
	n := new(AVLTNode[K, V])
	n.keyVal = KeyVal[K, V]{}
	n.keyVal.key = key
	n.keyVal.value = value
	n.height = 0
	n.sizeLeft = 0
	n.sizeRight = 0
	n.father = t._NIL
	n.left = t._NIL
	n.right = t._NIL
	n.prevInserted = t._NIL
	n.nextInserted = t._NIL
	return n
}

func (t *AVLTree[K, V]) pushToLastInserted(n *AVLTNode[K, V]) {
	if t.root == t._NIL {
		t.firstInserted = n
		// self-link
		n.nextInserted = n
		n.prevInserted = n
		return
	}

	n.prevInserted = t.firstInserted.prevInserted
	n.nextInserted = t.firstInserted
	t.firstInserted.prevInserted.nextInserted = n
	t.firstInserted.prevInserted = n
}

func (t *AVLTree[K, V]) removeToLastInserted(n *AVLTNode[K, V]) {
	if t.root == t._NIL {
		return
	}
	if n == t.firstInserted {
		t.firstInserted = t.firstInserted.nextInserted
	}

	n.prevInserted.nextInserted = n.nextInserted
	n.nextInserted.prevInserted = n.prevInserted
}
func (t *AVLTree[K, V]) cleanNil() {
	t._NIL.father = t._NIL
	t._NIL.left = t._NIL
	t._NIL.right = t._NIL
	t._NIL.prevInOrder = t._NIL
	t._NIL.nextInOrder = t._NIL
	t._NIL.prevInserted = t._NIL
	t._NIL.nextInserted = t._NIL
	t._NIL.height = DEPTH_INITIAL
	t._NIL.sizeLeft = 0
	t._NIL.sizeRight = 0
}

func (t *AVLTree[K, V]) getNode(k K) *AVLTNode[K, V] {
	notFound := true
	n := t.root
	c := int64(0)
	for notFound && n != t._NIL {
		c = t.avlTreeConstructorParams.Comparator(k, n.keyVal.key)
		notFound = c != 0
		if notFound {
			if c > 0 {
				n = n.right
			} else {
				n = n.left
			}
		}
	}
	return n
}

func (t *AVLTree[K, V]) put(n *AVLTNode[K, V]) (V, error) {
	if t.IsEmpty() {
		t.size = 1
		t.root = n
		t.minValue = n
		t.cleanNil()
		// self linking
		n.nextInOrder = n
		n.prevInOrder = n

		// tracking the chronological order
		t.firstInserted = n
		// self linking
		n.nextInserted = n
		n.prevInserted = n
		return t.avlTreeConstructorParams.ValueZeroValue, nil
	}

	k := n.keyVal.key
	v := n.keyVal.value

	//x is the iterator, next is the next node to move on
	next := t.root // must not be set to NIL, due to the while condition
	x := t.root
	c := int64(0)
	// descend the tree
	stillSearching := true
	for stillSearching && (next != t._NIL) {
		x = next
		c = t.avlTreeConstructorParams.Comparator(k, x.keyVal.key)
		if c == 0 {
			if t.avlTreeConstructorParams.KeyCollisionBehavior == Replace {
				stillSearching = false
				oldValue := x.keyVal.value
				x.keyVal.key = k
				x.keyVal.value = v
				// since the node has been modified,
				t.removeToLastInserted(x)
				t.pushToLastInserted(x)
				return oldValue, nil
			} else if t.avlTreeConstructorParams.KeyCollisionBehavior == IgnoreInsertion {
				return x.keyVal.value, nil
			} else {
				// if (behavior == BehaviorOnKeyCollision.AddItsNotASet) // -> add
				// stillSearching = false
				// c = -1
				return t.avlTreeConstructorParams.ValueZeroValue,
					fmt.Errorf("unexpected key collision behaviour: %b", t.avlTreeConstructorParams.KeyCollisionBehavior)
			}
		}
		if stillSearching {
			if c > 0 {
				next = x.right
			} else {
				next = x.left
			}
		}
	}

	//the new node has been placed as a leaf, check for errors
	if next == t._NIL {
		// end of tree reached: x is a leaf
		if c > 0 {
			x.right = n
		} else {
			x.left = n
		}
		n.father = x
	} else {
		return t.avlTreeConstructorParams.ValueZeroValue, errors.New("NOT A END?")
	}
	if t.size != INTEGER_MAX_VALUE {
		t.size++
	}

	// adjust links for iterators
	if c > 0 {
		n.prevInOrder = x
		x.nextInOrder.prevInOrder = n
		n.nextInOrder = x.nextInOrder
		x.nextInOrder = n
	} else {
		n.nextInOrder = x
		n.prevInOrder = x.prevInOrder
		x.prevInOrder.nextInOrder = n
		x.prevInOrder = n
		if x == t.minValue {
			t.minValue = n // in the end
		}
	}

	// we don't use n: it's height is 0 and it's connected only to NIL -> is balanced
	t.insertFixup(x)
	t.cleanNil()

	// track chronological insertion
	t.pushToLastInserted(n)
	t._NIL.nextInserted = t._NIL
	t._NIL.prevInserted = t._NIL
	return v, nil
}

func (t *AVLTree[K, V]) insertFixup(n *AVLTNode[K, V]) {
	hl := int64(0)
	hr := int64(0)
	delta := int64(0)
	temp := n

	for n != t._NIL {

		hl = n.left.height
		hr = n.right.height
		if hl > hr { // max
			n.height = hl
		} else {
			n.height = hr
		}
		n.height++

		delta = hl - hr
		// adjust sizes
		temp = n.left
		if temp == t._NIL {
			hl = 0
		} else {
			hl = (temp.sizeLeft + temp.sizeRight + 1)
		}
		n.sizeLeft = hl
		temp = n.right
		if temp == t._NIL {
			hr = 0
		} else {
			hr = (temp.sizeLeft + temp.sizeRight + 1)
		}
		n.sizeRight = hr
		// update father's size, could be usefull
		temp = n.father
		if temp != t._NIL {
			if temp.left == n {
				temp.sizeLeft = (hl + hr + 1)
			} else {
				temp.sizeRight = (hl + hr + 1)
			}
		}
		if delta < 2 && delta > -2 {
			// no rotation
			n = n.father
		} else {
			t.rotate(delta >= 2, n)
			n = n.father.father
		}
	}
	t.cleanNil()
}

func (t *AVLTree[K, V]) rotate(isRight bool, n *AVLTNode[K, V]) {
	hl := int64(0)
	hr := int64(0)

	nSide := n // dummy assignment
	oldFather := n.father
	if isRight {
		// right
		nSide = n.left
		if nSide.right.height > nSide.left.height {
			// three-rotation : ignoring this difference would cause the tree to be
			// umbalanced again
			a := n
			b := nSide
			c := b.right
			if oldFather.left == a {
				oldFather.left = c
			} else {
				oldFather.right = c
			}
			c.father = oldFather
			//					}
			a.father = c
			a.left = c.right
			if c.right != t._NIL {
				c.right.father = a
			}
			c.right = a
			if c.left.father != t._NIL {
				c.left.father = b
			}
			b.right = c.left
			b.father = c
			c.left = b

			// recompute height
			c.height++
			a.height -= 2
			b.height--
			t._NIL.left = t._NIL
			t._NIL.right = t._NIL
			t._NIL.father = t._NIL
			if a == t.root {
				t.root = c
				//				c.father = t._NIL // not necessary, but done to be sure
			}
			a.sizeLeft = c.sizeRight
			b.sizeRight = c.sizeLeft
			c.sizeRight += 1 + a.sizeRight
			c.sizeLeft += 1 + b.sizeLeft
			return
		}
		n.left = n.left.right // i could have put "nSide. .." but the whole piece of code would be less clear

		nSide.right.father = n
		nSide.right = n
		// adjust sizes
		n.sizeLeft = nSide.sizeRight
		nSide.sizeRight += 1 + n.sizeRight
	} else {
		// left
		nSide = n.right
		if nSide.left.height > nSide.right.height {
			a := n
			b := nSide
			c := b.left
			if oldFather.left == a {
				oldFather.left = c
			} else {
				oldFather.right = c
			}
			c.father = oldFather
			a.father = c
			a.right = c.left
			if c.left != t._NIL {
				c.left.father = a
			}
			c.left = a
			if c.right.father != t._NIL {
				c.right.father = b
			}
			b.left = c.right
			b.father = c
			c.right = b

			// recompute height
			c.height++
			a.height -= 2
			b.height--
			t._NIL.left = t._NIL
			t._NIL.right = t._NIL
			t._NIL.father = t._NIL
			if a == t.root {
				t.root = c
				//						c.father = NIL // not necessary, but done to be sure
			}
			a.sizeRight = c.sizeLeft
			b.sizeLeft = c.sizeRight
			c.sizeLeft += 1 + a.sizeLeft
			c.sizeRight += 1 + b.sizeRight
			return
		}
		n.right = n.right.left
		nSide.left.father = n
		nSide.left = n
		// adjust sizes
		n.sizeRight = nSide.sizeLeft
		nSide.sizeLeft += 1 + n.sizeLeft
	}
	if t._NIL == (oldFather) {
		// i'm root .. ehm, no: i WAS root
		t.root = nSide
	} else {
		if oldFather.left == n {
			oldFather.left = nSide
		} else {
			oldFather.right = nSide
		}
	}
	n.father = nSide
	nSide.father = oldFather
	hl = n.left.height
	hr = n.right.height
	if hl > hr {
		n.height = hl + 1
	} else {
		n.height = hr + 1
	}
	hl = n.height
	if isRight {
		hr = nSide.left.height
	} else {
		hr = nSide.right.height
	}
	if hl > hr {
		nSide.height = hl + 1
	} else {
		nSide.height = hr + 1
	}
}

func (t *AVLTree[K, V]) remove(n *AVLTNode[K, V]) (V, error) {
	if n == nil || t.IsEmpty() || n == t._NIL {
		return t.avlTreeConstructorParams.ValueZeroValue, nil // still might get a null-pointer-error
	}

	v := n.keyVal.value
	if t.size == 1 {
		if t.root != n && t.avlTreeConstructorParams.Comparator(t.root.keyVal.key, n.keyVal.key) == 0 {
			v = t.root.keyVal.value
		}
		t.size = 0
		t.root = t._NIL
		//cleanup
		t.minValue = t._NIL
		t.firstInserted = t._NIL
		t.cleanNil()
		return v, nil
	}

	prevSize := t.size
	// the real deletion
	ntbdFather := n.father
	actionPosition := ntbdFather
	notNilFather := ntbdFather != t._NIL
	hasLeft := n.left != t._NIL
	hasRight := n.right != t._NIL

	fmt.Printf("\n\n ON REMOVING -> current node: %v\n\thas left: %t, has right: %t\n", n.keyVal.key, hasLeft, hasRight)
	if hasLeft && hasRight {
		successor := n.nextInOrder // 33

		fmt.Println("has both left and right")
		fmt.Printf("successor: < k: %v> , isLeaf: %t\n", successor.keyVal.key, t.isLeaf(successor))
		fmt.Printf("successor father: < k: %v>\n", successor.father.keyVal.key)

		// shift values
		n.keyVal.key = successor.keyVal.key
		n.keyVal.value = successor.keyVal.value
		actionPosition = successor

		// but first, prepare non-inOrder optimization's links
		nPrevInserted := n.prevInserted
		nNextInserted := n.nextInserted
		successorSPrevInserted := successor.prevInserted
		successorSNextInserted := successor.nextInserted

		t.remove(successor)
		//t.updateOptimizationsOnRemove(n, successor)
		// un-link "n"
		nPrevInserted.nextInserted = nNextInserted
		nNextInserted.prevInserted = nPrevInserted
		// re-link successors' neighbor
		successorSPrevInserted.nextInserted = n
		n.prevInserted = successorSPrevInserted
		successorSNextInserted.prevInserted = n
		n.nextInserted = successorSNextInserted

		t.recalculateHeight(successor, true)
		t.recalculateSizes(successor, true)

	} else if hasLeft || hasRight {
		// just one child -> that child is a leaf
		// otherwise, a rotation would have happened while insertion
		// causing this node to become a fully-branched node (2 children)

		var child *AVLTNode[K, V]
		if hasLeft {
			child = n.left
			n.left = t._NIL
		} else {
			child = n.right
			n.right = t._NIL
		}
		// just shilft values
		n.keyVal.key = child.keyVal.key
		n.keyVal.value = child.keyVal.value
		n.height = 0
		n.sizeLeft = 0
		n.sizeRight = 0
		child.father = t._NIL
		actionPosition = n

		t.updateOptimizationsOnRemove(n, child)
	} else {
		// leaf -> nodeToUnlink = n
		if notNilFather {
			if ntbdFather.left == n {
				ntbdFather.left = t._NIL
				ntbdFather.sizeLeft = 0
			} else {
				ntbdFather.right = t._NIL
				ntbdFather.sizeRight = 0
			}
			// adjust links
			t.unlinkUpdateOptimizations(n)
		} else {
			// n is both root and leaf -> empty
			t.root = t._NIL
			t.size = 0
			// other clean-ups will be performed later
			// now the tree is empty
			actionPosition = t._NIL
		}
	}

	t.cleanNil()

	//then, balance
	if actionPosition != t._NIL {
		t.recalculateHeight(actionPosition, true)
		t.recalculateSizes(actionPosition, true)
		t.insertFixup(actionPosition)
	}

	t.size--
	if t.size == 0 || t.IsEmpty() {
		t.root = t._NIL
		t.size = 0
		t.minValue = t._NIL
		t.firstInserted = t._NIL
	} else if t.root == t._NIL {
		return t.avlTreeConstructorParams.ValueZeroValue, fmt.Errorf("BUG: root should not be nil since size is not zero")
	}

	// adjusting connections
	if t.size == 1 {
		t.root.father = t._NIL
		t.minValue = t.root
		t.root.nextInOrder = t.root
		t.root.prevInOrder = t.root
		t.firstInserted = t.root
		t.root.nextInserted = t.root
		t.root.prevInserted = t.root
	}

	if prevSize == INTEGER_MAX_VALUE {
		prevSize = 1 + t.root.sizeLeft + t.root.sizeRight
		if prevSize < 0 || t.root.sizeLeft < 0 || t.root.sizeRight < 0 {
			prevSize = INTEGER_MAX_VALUE
		}
		t.size = prevSize
	}
	t.cleanNil()

	return v, nil
}

func (t *AVLTree[K, V]) index(n *AVLTNode[K, V]) int {
	i := 0
	if n.sizeLeft > 0 {
		i = int(n.sizeLeft)
	}
	indexNodeLeftInorder := -1
	nodeLeftInOrder := n
	// travel back the "left full-branching" until we find the right branching
	// this way, the father causing the right branching will be the "previous in order"
	// of the "n" node. Add its index to this "n" current one

	for nodeLeftInOrder != t._NIL && nodeLeftInOrder.father != t._NIL && nodeLeftInOrder.father.right != nodeLeftInOrder {
		nodeLeftInOrder = nodeLeftInOrder.father
	}
	if nodeLeftInOrder != t._NIL && nodeLeftInOrder.father != t._NIL { // && n.father.right == n
		indexNodeLeftInorder = t.index(nodeLeftInOrder.father)
	} //else { indexNodeLeftInorder = -1 }
	if indexNodeLeftInorder >= 0 {
		i += indexNodeLeftInorder + 1
	}
	return i
}

func (t *AVLTree[K, V]) isLeaf(n *AVLTNode[K, V]) bool {
	return t == nil || n == nil || (n.left == t._NIL && n.right == t._NIL)
}

func (t *AVLTree[K, V]) recalculateHeight(n *AVLTNode[K, V], recurseToRoot bool) {
	if t._NIL == n {
		return
	}
	var h int64
	shouldContinue := true
	for shouldContinue {
		h = DEPTH_INITIAL
		if n.left != t._NIL {
			h = n.left.height
		}
		if n.right != t._NIL && n.right.height > h {
			h = n.right.height
		}
		n.height = 1 + h

		if recurseToRoot {
			n = n.father
			shouldContinue = (n != t._NIL)
		} else {
			shouldContinue = false
		}

	}
}
func (t *AVLTree[K, V]) recalculateSizes(n *AVLTNode[K, V], recurseToRoot bool) {
	if t._NIL == n {
		return
	}
	shouldContinue := true
	for shouldContinue {
		if n.left != t._NIL {
			n.sizeLeft = 1 + n.left.sizeLeft + n.left.sizeRight
		} else {
			n.sizeLeft = 0
		}
		if n.right != t._NIL {
			n.sizeRight = 1 + n.right.sizeLeft + n.right.sizeRight
		} else {
			n.sizeRight = 0
		}

		if recurseToRoot {
			n = n.father
			shouldContinue = (n != t._NIL)
		} else {
			shouldContinue = false
		}

	}
}

func (n *AVLTNode[K, V]) unlinkAll() {
	n.prevInOrder.nextInOrder = n.nextInOrder
	n.nextInOrder.prevInOrder = n.prevInOrder

	n.prevInserted.nextInserted = n.nextInserted
	n.nextInserted.prevInserted = n.prevInserted
}

func (t *AVLTree[K, V]) unlinkUpdateOptimizations(n *AVLTNode[K, V]) {
	if t.minValue == n {
		t.minValue = n.nextInOrder
	}
	if t.firstInserted == n {
		t.firstInserted = n.nextInserted
	}
	n.unlinkAll()
}

func (t *AVLTree[K, V]) updateOptimizationsOnRemove(nWillBeSwapped *AVLTNode[K, V], childWillBeDestroyed *AVLTNode[K, V]) {
	// NOTE: I'll leave the internal comments just to reference
	// and explain the thought processes

	// CHRONOLOGICAL ORDERING (FIFO)

	if t.firstInserted == nWillBeSwapped {
		t.firstInserted = nWillBeSwapped.nextInserted

		// A.1)
		// the "n" node should be forgotten since it will be removed
		// so, at first, let unlink it

		// A.2)
		// since "child" will take "N"'s place, shift the links towards
		// that node instance ...

		// A.3)
		// ... and the "new-child"'s links towards the previous "chronological-neighbour"
	} else if t.firstInserted == childWillBeDestroyed {
		t.firstInserted = nWillBeSwapped

		// B.1)
		// those who were pointing to "child" still need to point to the
		// same <key; value> pair, so redirect them to "n"

		// B.2)
		// those who were pointing to "n" needs to "forget" it
		// since it will be removed

		// B.3)
		// the "child" <key, value> pair needs to preserve the insertion
		// order, so let's keep it
	}
	// A.1 & B.2
	nWillBeSwapped.prevInserted.nextInserted = nWillBeSwapped.nextInserted
	nWillBeSwapped.nextInserted.prevInserted = nWillBeSwapped.prevInserted

	// A.2 & B.1
	childWillBeDestroyed.prevInserted.nextInserted = nWillBeSwapped
	childWillBeDestroyed.nextInserted.prevInserted = nWillBeSwapped

	// A.3 & B.3
	nWillBeSwapped.prevInserted = childWillBeDestroyed.prevInserted
	nWillBeSwapped.nextInserted = childWillBeDestroyed.nextInserted

	// cleanse the "old child"
	childWillBeDestroyed.prevInserted = t._NIL
	childWillBeDestroyed.nextInserted = t._NIL

	// MIN-VALUE

	if t.minValue == childWillBeDestroyed {
		// redirect links towards "n", since the "min key" has to remain the same
		t.minValue = nWillBeSwapped
	} else if t.minValue == nWillBeSwapped {
		// just update the value
		t.minValue = nWillBeSwapped.nextInOrder
	}
	/*
		// "n"'s neighbour needs to forget that node
		nWillBeSwapped.prevInOrder.nextInOrder = nWillBeSwapped.nextInOrder
		nWillBeSwapped.nextInOrder.prevInOrder = nWillBeSwapped.prevInOrder

		// the "new child" needs to remember its previous neighbours
		nWillBeSwapped.prevInOrder = childWillBeDestroyed.prevInOrder
		nWillBeSwapped.nextInOrder = childWillBeDestroyed.nextInOrder

		// the "old child" neighbours now need to track the right node: "n"
		childWillBeDestroyed.nextInOrder.prevInOrder = nWillBeSwapped
		childWillBeDestroyed.prevInOrder.nextInOrder = nWillBeSwapped
	*/
	//unlink the "child"
	childWillBeDestroyed.nextInOrder.prevInOrder = childWillBeDestroyed.prevInOrder
	childWillBeDestroyed.prevInOrder.nextInOrder = childWillBeDestroyed.nextInOrder

	// the "Old n" node instance will hold the "child"'s data,
	// so on need to update links

	// cleanse the "old child"
	childWillBeDestroyed.prevInOrder = t._NIL
	childWillBeDestroyed.nextInOrder = t._NIL
}

func (t *AVLTree[K, V]) toStringTabbed(fullLogNode bool, n *AVLTNode[K, V], tabLevel int, printer func(string)) {
	var sb strings.Builder
	for i := 0; i < tabLevel; i++ {
		sb.WriteString("  ")
	}
	sb.WriteString("-> ")
	if n == nil || n == t._NIL {
		sb.WriteString("null")
		printer(sb.String())
		return
	}

	sb.WriteString(" -- at index= ")
	sb.WriteString(strconv.Itoa(t.index(n)))
	sb.WriteString(" we have ")
	n.toStringTabbed(fullLogNode, func(s string) { sb.WriteString(s) })
	sb.WriteString(" ;; and value= <<")
	val, ok := any(n.keyVal.value).(*AVLTree[K, V])
	if ok && val == t {
		sb.WriteString(" SELF TREE - RECURSION AVOIDED ")
	} else {
		sb.WriteString(fmt.Sprint(n.keyVal.value))
	}
	sb.WriteString("\n")
	printer(sb.String())
	sb.Reset() // make it invalid / unusable
	t.toStringTabbed(fullLogNode, n.left, tabLevel+1, printer)
	printer("\n")
	t.toStringTabbed(fullLogNode, n.right, tabLevel+1, printer)
	printer("\n")
}

func (n *AVLTNode[K, V]) toStringTabbed(fullLogNode bool, printer func(string)) {
	var sb strings.Builder
	sb.WriteString("Node [ ")
	sb.WriteString("key= <<")
	sb.WriteString(fmt.Sprint(n.keyVal.key))
	sb.WriteString(">>")
	if fullLogNode {
		sb.WriteString(" ; height= ")
		sb.WriteString(strconv.FormatInt(n.height, 10))
		sb.WriteString(" ;; size left= ")
		sb.WriteString(strconv.FormatInt(n.sizeLeft, 10))
		sb.WriteString(" ; size right= ")
		sb.WriteString(strconv.FormatInt(n.sizeRight, 10))
		sb.WriteString(" ;; father's key= <<")
		sb.WriteString(fmt.Sprint(n.father.keyVal.key))
		sb.WriteString(">> ;; next-in-order's key= <<")
		sb.WriteString(fmt.Sprint(n.nextInOrder.keyVal.key))
		sb.WriteString(">> ;; prev-in-order's key= <<")
		sb.WriteString(fmt.Sprint(n.prevInOrder.keyVal.key))
		sb.WriteString(">> ;; next-chronological's key= <<")
		sb.WriteString(fmt.Sprint(n.nextInserted.keyVal.key))
		sb.WriteString(">> ;; prev-chronological's key= <<")
		sb.WriteString(fmt.Sprint(n.prevInserted.keyVal.key))
		sb.WriteString(">>")
	}
	sb.WriteString(" ]")
	printer(sb.String())
}

func (t *AVLTree[K, V]) inlineStringKeyOnly(printer func(string), n *AVLTNode[K, V]) {
	printer("(")
	if n != t._NIL && n != nil {
		t.inlineStringKeyOnly(printer, n.left)
		printer(" -> ")
		printer(fmt.Sprintf("<h:%d; k:%v>", n.height, n.keyVal.key))
		printer(" <- ")
		t.inlineStringKeyOnly(printer, n.right)
	}
	printer(")")
}

func (t *AVLTree[K, V]) forEachNode(mode ForEachMode, action func(*AVLTNode[K, V])) error {
	if t.IsEmpty() || action == nil {
		return nil
	}
	canContinue := true
	iterMax := t.size + 1 //anti-bug
	switch mode {
	case InOrder:
		{
			current := t.minValue
			start := current
			for canContinue && iterMax >= 0 { // do-while loop
				iterMax--
				if iterMax < 0 {
					return fmt.Errorf("BUG ! for-each is looping more than expected")
				}
				action(current)
				current = current.nextInOrder
				canContinue = current != start
			}
		}
	case ReverseInOrder:
		{
			current := t.minValue.prevInOrder
			start := current
			for canContinue && iterMax >= 0 { // do-while loop
				iterMax--
				if iterMax < 0 {
					return fmt.Errorf("BUG ! for-each is looping more than expected")
				}
				action(current)
				current = current.prevInOrder
				canContinue = current != start
			}
		}
	case Stack:
		{
			current := t.firstInserted.prevInserted
			start := current
			for canContinue && iterMax >= 0 { // do-while loop
				iterMax--
				if iterMax < 0 {
					return fmt.Errorf("BUG ! for-each is looping more than expected")
				}
				action(current)
				current = current.prevInserted
				canContinue = current != start
			}
		}
	case Queue:
		{
			current := t.firstInserted
			start := current
			for canContinue && iterMax >= 0 { // do-while loop
				iterMax--
				if iterMax < 0 {
					return fmt.Errorf("BUG ! for-each is looping more than expected")
				}
				action(current)
				current = current.nextInserted
				canContinue = current != start
			}
		}
	}
	return nil
}

//

// PUBLIC functions

//

func NewAVLTree[K any, V any](avlTreeConstructorParams AVLTreeConstructorParams[K, V]) (*AVLTree[K, V], error) {
	if avlTreeConstructorParams.KeyExtractor == nil {
		return nil, errors.New("key extractor must not be null")
	}
	if avlTreeConstructorParams.Comparator == nil {
		return nil, errors.New("comparator must not be null")
	}
	t := new(AVLTree[K, V])
	t.avlTreeConstructorParams = avlTreeConstructorParams
	t.size = 0
	t._NIL = nil
	_nil := t.newNode(avlTreeConstructorParams.KeyZeroValue, avlTreeConstructorParams.ValueZeroValue)
	t._NIL = _nil
	t.root = _nil
	t.cleanNil()

	// other optimization-based setup
	t.firstInserted = _nil
	t.minValue = _nil

	return t, nil
}

func (t *AVLTree[K, V]) Size() int64 {
	return t.size
}

func (t *AVLTree[K, V]) NILL() interface{} { return t._NIL }

func (t *AVLTree[K, V]) Put(key K, value V) (V, error) {
	n := t.newNode(key, value)
	return t.put(n)
}

func (t *AVLTree[K, V]) Remove(k K) (V, error) {
	if t.root == t._NIL {
		return t.avlTreeConstructorParams.ValueZeroValue, nil
	}
	n := t.getNode(k)
	if n == t._NIL {
		return t.avlTreeConstructorParams.ValueZeroValue, fmt.Errorf("key not found")
	}
	return t.remove(n)
}

func (t *AVLTree[K, V]) IsEmpty() bool {
	return t == nil || t.root == t._NIL
}

func (t *AVLTree[K, V]) ForEach(mode ForEachMode, action func(K, V)) error {
	if t.IsEmpty() || action == nil {
		return nil
	}
	return t.forEachNode(mode, func(n *AVLTNode[K, V]) { action(n.keyVal.key, n.keyVal.value) })
}

func (t *AVLTree[K, V]) StringInto(fullLogNode bool, printer func(string)) {
	printer("AVL Tree ")
	if t == nil {
		printer("\t- NULL!")
		return
	}
	if t.root == t._NIL {
		printer("\t- empty")
	}
	printer("of size= ")
	printer(strconv.FormatInt(t.size, 10))
	printer("; :\n")
	t.toStringTabbed(fullLogNode, t.root, 0, printer)
}
func (t *AVLTree[K, V]) StringLogginFull(fullLogNode bool) string {
	var sb strings.Builder
	t.StringInto(fullLogNode, func(s string) { sb.WriteString(s) })
	return sb.String()
}
func (t *AVLTree[K, V]) String() string {
	return t.StringLogginFull(true)
}
func (t *AVLTree[K, V]) InlineStringKeyOnly(printer func(string)) {
	t.inlineStringKeyOnly(printer, t.root)
}

func (n *AVLTNode[K, V]) String() string {
	var sb strings.Builder
	n.toStringTabbed(true, func(s string) { sb.WriteString(s) })
	return sb.String()
}

//

func (p *KeyVal[K, V]) PairKey() K {
	return p.key
}
func (p *KeyVal[K, V]) PairValue() V {
	return p.value
}

func (b ForEachMode) String() string {
	var name string
	switch b {
	case InOrder:
		name = "InOrder"
	case ReverseInOrder:
		name = "ReverseInOrder"
	case Queue:
		name = "Queue"
	case Stack:
		name = "Stack"
	default:
		name = fmt.Sprintf("unrecognized ForEachMode: %b", b)
	}
	return name
}
