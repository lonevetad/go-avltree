package avltree

import (
	"errors"
	"fmt"
	"strconv"
	"strings"
)

const DEPTH_NIL int64 = -1
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

func (t *AVLTree[K, V]) cleanNode(n *AVLTNode[K, V]) {
	n.father = t._NIL
	n.left = t._NIL
	n.right = t._NIL
	n.prevInOrder = t._NIL
	n.nextInOrder = t._NIL
	n.prevInserted = t._NIL
	n.nextInserted = t._NIL
	n.height = 0
	n.sizeLeft = 0
	n.sizeRight = 0
}

func (t *AVLTree[K, V]) cleanNil() {
	t.cleanNode(t._NIL)
	t._NIL.height = DEPTH_NIL
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
		n.height = 0
		n.sizeLeft = 0
		n.sizeRight = 0
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
	return t.avlTreeConstructorParams.ValueZeroValue, nil
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
			t.rotate(delta, n)
			n = n.father.father
		}
	}
	t.cleanNil()
}

func (t *AVLTree[K, V]) rotate(delta int64, n *AVLTNode[K, V]) {
	isRight := delta >= 2
	hl := int64(0)
	hr := int64(0)

	nSide := n // dummy assignment
	oldFather := n.father
	if isRight {
		// right
		nSide = n.left
		if nSide.right.height > nSide.left.height {
			// left-right rotation (left on b, right on c)
			// three-rotation : ignoring this difference would cause the tree to be
			// umbalanced again. x,y,z and w might be (all or none) NIL
			//h  :   .   n=a
			//   .   ./  .  \.
			//h+1:  b.   .   x
			//h+2:y  .c
			//h+3:   z w
			// ->
			//h  :   c
			//   .  /.\
			//h+1: b . a
			//h+2:y z.w x.
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
			if c.left.father != t._NIL { // c.left.father == c, but this would modify "NIL" if "c == NIL"
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
				c.father = t._NIL // not necessary, but done to be sure
			}
			a.sizeLeft = c.sizeRight
			b.sizeRight = c.sizeLeft
			c.sizeRight += 1 + a.sizeRight
			c.sizeLeft += 1 + b.sizeLeft
			return
		}
		// right rotation on b
		//h  :   .   n=a
		//   .   ./  .  \.
		//h+1:  b.   .   x
		//h+2: c .y
		//h+3:z w
		// ->
		//h  :   b
		//   .  /.\
		//h+1: c . a
		//h+2:z w.y x.

		// note: oldFather pointers moved outside
		n.left = n.left.right // a.left == y
		nSide.right.father = n
		nSide.right = n // b.right == a

		// adjust sizes
		n.sizeLeft = nSide.sizeRight
		nSide.sizeRight += 1 + n.sizeRight
	} else {
		// left
		nSide = n.right
		if nSide.left.height > nSide.right.height {
			// right-left rotation (right on b, left on c)
			// three-rotation : ignoring this difference would cause the tree to be
			// umbalanced again. x,y,z and w might be (all or none) NIL
			//h  :   .   n=a
			//   .   ./  .  \.
			//h+1:  x.   .   b
			//h+2:   .   .c  .  y
			//h+3:   .   z w
			// ->
			//h  :   c
			//   .  /.\
			//h+1: a . b
			//h+2:x z.w y.
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
				c.father = t._NIL // not necessary, but done to be sure
			}
			a.sizeRight = c.sizeLeft
			b.sizeLeft = c.sizeRight
			c.sizeLeft += 1 + a.sizeLeft
			c.sizeRight += 1 + b.sizeRight
			return
		}
		// left rotation on b
		//h  :   .   n=a
		//   .   ./  .  \.
		//h+1:  x.   .   b
		//h+2:   .   . y . c
		//h+3:   .   .   .z w
		// ->
		//h  :   b
		//   .  /.\
		//h+1: a . c
		//h+2:x y.z w.

		// note: oldFather pointers moved outside
		n.right = n.right.left // y
		nSide.left.father = n
		nSide.left = n
		// adjust sizes
		n.sizeRight = nSide.sizeLeft
		nSide.sizeLeft += 1 + n.sizeLeft
	}
	// other clean-ups, adjustments, etc
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
	if hl > hr { // isRight
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

	if hasLeft && hasRight {
		successor := n.nextInOrder // 33

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

func (t *AVLTree[K, V]) index(n *AVLTNode[K, V]) int64 {
	if n == t._NIL {
		return 0
	}
	i := n.sizeLeft

	lowestGrandDadOnLeft := n
	// are we the root? That means, is there no other node
	// somewhere on the "left" (in a lower height)?
	for lowestGrandDadOnLeft != t._NIL && lowestGrandDadOnLeft == lowestGrandDadOnLeft.father.left {
		lowestGrandDadOnLeft = lowestGrandDadOnLeft.father
	}
	//if exists a node "M" having "n" as right, then "n" is on the
	// right side of a bigger sub-tree, whose root is "M". Get M's
	// index, then add n's left part (leftSize) then count n
	// itself (+1) in order to include M in the "total noudes count"
	if lowestGrandDadOnLeft != t._NIL && lowestGrandDadOnLeft.father != t._NIL && (lowestGrandDadOnLeft.father.right == lowestGrandDadOnLeft) {
		i += 1 + t.index(lowestGrandDadOnLeft.father)
	}
	return i
}

func (t *AVLTree[K, V]) isLeaf(n *AVLTNode[K, V]) bool {
	return t == nil || n == nil || (n.left == t._NIL && n.right == t._NIL)
}

func (t *AVLTree[K, V]) getNodeAt(index int64) *AVLTNode[K, V] {
	if t.IsEmpty() || t.size < 1 || index < 0 || index >= t.size {
		return nil
	}
	if t.size == 1 {
		return t.root
	}
	n := t.root
	var nprev *AVLTNode[K, V]
	nprev = t.root

	notFound := true
	for notFound && index >= 0 && n != t._NIL {
		nprev = n

		if index < n.sizeLeft {
			n = n.left
		} else {
			index -= n.sizeLeft
			if index > 0 {
				index-- // remove current node from the "count" of the "available one"
				n = n.right
			} else {
				notFound = false // FOUND! close the loop
			}
		}
	}
	return nprev
}

/*
the assignments to "left" and "right" pointers disrupt the tree, whose integrity is
fundamental to make "getNodeAt" properly work.
So, the "left" nodes are placed in the "prev in order" pointers while the "right"
nodes are placed in the "next in order" pointers.
*/
func (t *AVLTree[K, V]) getRootForSubtreeOnCompacting(lowerIndex, upperIndex int64, dept int) *AVLTNode[K, V] {
	var newLeft, newRight *AVLTNode[K, V]

	if upperIndex <= lowerIndex {
		currentNode := t.getNodeAt(upperIndex)
		currentNode.height = 0
		// mark me as a leaf
		currentNode.prevInOrder = t._NIL
		currentNode.nextInOrder = t._NIL
		currentNode.father = t._NIL
		return currentNode
	}

	currentIndex := (lowerIndex + upperIndex) >> 1

	currentNode := t.getNodeAt(currentIndex)
	for i := dept; i > 0; i-- {
		fmt.Print("\t")
	}
	fmt.Printf("at dept %d, currentIndex %d, found node < %v >\n", dept, currentIndex, currentNode.keyVal.key)

	if lowerIndex < currentIndex {
		newLeft = t.getRootForSubtreeOnCompacting(lowerIndex, currentIndex-1, dept+1)
	} else {
		newLeft = t._NIL
	}
	if currentIndex < upperIndex {
		newRight = t.getRootForSubtreeOnCompacting(currentIndex+1, upperIndex, dept+1)
	} else {
		newRight = t._NIL
	}
	t.cleanNil()

	if newLeft == t._NIL {
		currentNode.prevInOrder = t._NIL
	} else {
		currentNode.prevInOrder = newLeft
	}
	if newRight == t._NIL {
		currentNode.nextInOrder = t._NIL
	} else {
		currentNode.nextInOrder = newRight
	}

	t.cleanNil()

	return currentNode
}

/* returns: the nodes holding the minimum and maximum keys respectively (among the nodes
* in this sub-tree)
 */
func (t *AVLTree[K, V]) reconstructTreeOnCompactingOnRoot(currentRoot *AVLTNode[K, V]) (*AVLTNode[K, V], *AVLTNode[K, V]) {
	if t._NIL == currentRoot {
		// kinda-error
		return t._NIL, t._NIL
	}
	if currentRoot.prevInOrder == t._NIL && currentRoot.nextInOrder == t._NIL {
		// I'm a LEAF!
		currentRoot.left = t._NIL
		currentRoot.right = t._NIL
		currentRoot.height = 0
		currentRoot.sizeLeft = 0
		currentRoot.sizeRight = 0
		return currentRoot, currentRoot // a leaf holds both the min and max values of a itself-tree
	}

	var minNode, maxNode, newPrev, newNext *AVLTNode[K, V]

	minNode, newPrev = t.reconstructTreeOnCompactingOnRoot(currentRoot.prevInOrder)
	newNext, maxNode = t.reconstructTreeOnCompactingOnRoot(currentRoot.nextInOrder)

	//re-link left, right and fathers, then sizes left/right and height
	currentRoot.left = currentRoot.prevInOrder
	if currentRoot.prevInOrder == t._NIL {
		currentRoot.sizeLeft = 0
	} else {
		currentRoot.left.father = currentRoot
		currentRoot.sizeLeft = currentRoot.left.sizeLeft + 1 + currentRoot.left.sizeRight
	}
	currentRoot.right = currentRoot.nextInOrder
	if currentRoot.nextInOrder == t._NIL {
		currentRoot.sizeRight = 0
	} else {
		currentRoot.right.father = currentRoot
		currentRoot.sizeRight = currentRoot.right.sizeLeft + 1 + currentRoot.right.sizeRight
	}
	//height
	if currentRoot.left.height > currentRoot.right.height {
		currentRoot.height = currentRoot.left.height + 1
	} else {
		currentRoot.height = currentRoot.right.height + 1
	}

	// re-link prev/next in order
	if newPrev != t._NIL {
		newPrev.nextInOrder = currentRoot
		currentRoot.prevInOrder = newPrev
	}
	if newNext != t._NIL {
		newNext.prevInOrder = currentRoot
		currentRoot.nextInOrder = newNext
	}

	if minNode == t._NIL {
		minNode = currentRoot
	}
	if maxNode == t._NIL {
		maxNode = currentRoot
	}
	return minNode, maxNode
}

func (t *AVLTree[K, V]) reconstructTreeOnCompacting(currentRoot *AVLTNode[K, V]) {
	minNode, maxNode := t.reconstructTreeOnCompactingOnRoot(currentRoot)
	// close the loop
	minNode.prevInOrder = maxNode
	maxNode.nextInOrder = minNode
}

func (t *AVLTree[K, V]) recalculateHeight(n *AVLTNode[K, V], recurseToRoot bool) {
	if t._NIL == n {
		return
	}
	var h int64
	shouldContinue := true
	for shouldContinue {
		h = DEPTH_NIL
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
	sb.WriteString(strconv.FormatInt(t.index(n), 10))
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

func (t *AVLTree[K, V]) Clear() {
	t.size = 0
	t.root = t._NIL
	t.minValue = t._NIL
	t.firstInserted = t._NIL
	t.cleanNil()
}

func (t *AVLTree[K, V]) NILL() interface{} { return t._NIL }

func (t *AVLTree[K, V]) GetAt(index int64) (*KeyVal[K, V], error) {
	n := t.getNodeAt(index)
	if n == nil {
		return nil, fmt.Errorf("no node to retrieve")
	}
	return &n.keyVal, nil
}

func (t *AVLTree[K, V]) Put(key K, value V) (V, error) {
	n := t.newNode(key, value)
	return t.put(n)
}

func (t *AVLTree[K, V]) Remove(k K) (V, error) {
	if t.root == t._NIL {
		return t.avlTreeConstructorParams.ValueZeroValue, EMPTY_TREE()
	}
	n := t.getNode(k)
	if n == t._NIL {
		return t.avlTreeConstructorParams.ValueZeroValue, KEY_NOT_FOUND()
	}
	return t.remove(n)
}

func (t *AVLTree[K, V]) IsEmpty() bool {
	return t == nil || t.root == t._NIL
}

func (t *AVLTree[K, V]) CompactBalance() {
	if t.size < 8 {
		/* Adding a series of sorted numbers leads
		to a heavily "unbalanced" tree towards a single side,
		creating a fractal-like structure whose sizes
		follows the Fibonacci numbers.
		3 nodes lead to a obviously balanced tree.
		5 nodes lead to a tree that seems prone to umbalance,
		but it is not; meanwile a six-th node would cause
		a rotation to happen.
		A rotation involving the root happen after increasing the
		size over a Fibonacci number and the next one is 8.
		7 leads to a full tree, while 8 increase the tree height
		without a rotation. So 8 is the threshold.
		*/
		return
	}

	newRoot := t.getRootForSubtreeOnCompacting(0, t.size-1, 0)
	t.cleanNil()
	t.root = newRoot
	t.reconstructTreeOnCompacting(newRoot)
	fmt.Printf("new root (key: %v)\n", t.root.keyVal.key)
	fmt.Println("new tree:")
	fmt.Println(t)
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

func (p *KeyVal[K, V]) Key() K {
	return p.key
}
func (p *KeyVal[K, V]) Value() V {
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
