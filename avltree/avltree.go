package avltree

import (
	"errors"
	"fmt"
	"strconv"
	"strings"
)

const DEPTH_INITIAL int = -1
const INTEGER_MAX_VALUE int = 2147483647

type ForEachMode byte

const (
	InOrder        ForEachMode = 0
	ReverseInOrder ForEachMode = 1
	Queue          ForEachMode = 2
	Stack          ForEachMode = 3
)

// TODO: an enum for "for each" order"

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

type Comparator[K any] func(k1 K, k2 K) int

type AVLTNode[K any, V any] struct {
	keyVal       KeyVal[K, V]
	height       int
	sizeLeft     int
	sizeRight    int
	father       *AVLTNode[K, V]
	left         *AVLTNode[K, V]
	right        *AVLTNode[K, V]
	nextInOrder  *AVLTNode[K, V]
	prevInOrder  *AVLTNode[K, V]
	prevInserted *AVLTNode[K, V]
	nextInserted *AVLTNode[K, V]
}

type AVLTree[K any, V any] struct {
	size           int
	keyZeroValue   K
	valueZeroValue V
	keyExtractor   KeyExtractor[K, V]
	comparator     Comparator[K]
	root           *AVLTNode[K, V]
	_NIL           *AVLTNode[K, V]
	minValue       *AVLTNode[K, V] // used for optimizations
	firstInserted  *AVLTNode[K, V]
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
		t.firstInserted = t.firstInserted.prevInserted
	}

	n.prevInserted.nextInserted = n.nextInserted
	n.nextInserted.prevInserted = n.prevInserted
}

func (t *AVLTree[K, V]) put(n *AVLTNode[K, V]) (V, error) {
	prevSize := int(0)

	fmt.Printf("putting .... n: %v\n", n)
	if t.size == 0 || t.root == t._NIL {
		fmt.Print("FIRST PUT!")
		t.size = 1
		t.root = n
		t.minValue = n
		// self linking
		n.nextInOrder = n
		n.prevInOrder = n
		// cleaning "NIL"
		t._NIL.father = t._NIL
		t._NIL.left = t._NIL
		t._NIL.right = t._NIL
		t._NIL.nextInOrder = t._NIL
		t._NIL.prevInOrder = t._NIL

		//super.put(n)
		// tracking the chronological order
		t.firstInserted = n
		// self linking
		n.nextInserted = n
		n.prevInserted = n
		t._NIL.nextInserted = t._NIL
		t._NIL.prevInserted = t._NIL
		return t.valueZeroValue, nil
	}

	prevSize = t.size
	k := n.keyVal.key
	v := n.keyVal.value
	// v = super.put(n)

	//x is the iterator, next is the next node to move on
	next := t.root
	x := t.root // must not be set to NIL, due to the while condition
	// descend the tree
	c := int(0)
	stillSearching := true
	for stillSearching && (next != t._NIL) {
		x = next
		c = t.comparator(k, x.keyVal.key)
		if c == 0 {
			/*
				if (behaviour == MapTreeAVL.BehaviourOnKeyCollision.Replace) {
					oldValue = x.v
					x.k = k
					x.v = v
					return oldValue
				} else if (behaviour == MapTreeAVL.BehaviourOnKeyCollision.KeepPrevious)
					return x.v
				else // if (behavior == MapTreeAVL.BehaviorOnKeyCollision.AddItsNotASet) //
						// -> add
					c = -1
			*/
			stillSearching = false
			oldValue := x.keyVal.value
			x.keyVal.key = k
			x.keyVal.value = v

			fmt.Print("substituting ....")
			t.removeToLastInserted(n)
			t.pushToLastInserted(n)
			return oldValue, nil
		}
		if stillSearching {
			if c > 0 {
				next = x.right
			} else {
				next = x.left
			}
		}
	}

	if next == t._NIL {
		// end of tree reached: x is a leaf
		if c > 0 {
			x.right = n
		} else {
			x.left = n
		}
		n.father = x
	} else {
		return t.valueZeroValue, errors.New("NOT A END?")
	}
	if t.size != INTEGER_MAX_VALUE {
		t.size++
	}

	// adjust links for iterators
	// minValue
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

	// don't use n: it's height is 0 and it's connected only to NIL -> is balanced
	t.insertFixup(x)
	t._NIL.father = t._NIL
	t._NIL.left = t._NIL
	t._NIL.right = t._NIL
	t._NIL.nextInOrder = t._NIL
	t._NIL.prevInOrder = t._NIL

	if prevSize != INTEGER_MAX_VALUE && prevSize != t.size {
		// node really added
		n.prevInserted = t.firstInserted.prevInserted
		n.nextInserted = t.firstInserted
		t.firstInserted.prevInserted.nextInserted = n
		t.firstInserted.prevInserted = n
	}

	t._NIL.nextInserted = t._NIL
	t._NIL.prevInserted = t._NIL
	return v, nil
}

func (t *AVLTree[K, V]) insertFixup(n *AVLTNode[K, V]) {
	// TODO : paste ....	super.insertFixup(nnn)
	hl := int(0)
	hr := int(0)
	delta := int(0)
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
	t._NIL.sizeLeft = -1 // just an invalid value
	t._NIL.sizeRight = -1
	t._NIL.height = DEPTH_INITIAL

	t._NIL.nextInOrder = t._NIL
	t._NIL.prevInOrder = t._NIL
}

func (t *AVLTree[K, V]) rotate(isRight bool, n *AVLTNode[K, V]) {
	hl := int(0)
	hr := int(0)

	nSide := n // dummy assignment
	oldFather := n.father
	if isRight {
		// right
		nSide = n.left
		if nSide.right.height > nSide.left.height {
			// three-rotation : ignoring this difference would cause the tree to be
			// umbalanced again
			//final NodeAVL_Indexable a, b, c
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
			//final NodeAVL_Indexable a, b, c
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

/*
func (n *AVLTNode[K, V]) String() string {
	var sb strings.Builder
	sb.WriteString("k:")
	sb.WriteString(n.keyVal.key.String())
	sb.WriteString(" - v:")
	sb.WriteString(n.keyVal.value)
	sb.WriteString(",h:")
	sb.WriteString(n.height)
	sb.WriteString(",f:")
	sb.WriteString(n.father.keyVal.key)
	sb.WriteString(", sl:")
	sb.WriteString(n.sizeLeft)
	sb.WriteString(",sr:")
	sb.WriteString(n.sizeRight)
	return sb.String()
}
*/

func (t *AVLTree[K, V]) toStringTabbed(n *AVLTNode[K, V], tabLevel int, printer func(string)) {
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

	n.toStringTabbed(func(s string) { sb.WriteString(s) })
	sb.WriteString(" -- index= ")
	sb.WriteString(strconv.Itoa(t.index(n)))
	sb.WriteString(" ; father's key= <<")
	sb.WriteString(fmt.Sprint(n.father.keyVal.key))
	sb.WriteString(">> ;; value= <<")
	val, ok := any(n.keyVal.value).(*AVLTree[K, V])
	if ok && val == t {
		sb.WriteString(" SELF TREE - RECURSION AVOIDED ")
	} else {
		sb.WriteString(fmt.Sprint(n.keyVal.value))
	}
	sb.WriteString("\n")
	printer(sb.String())
	sb.Reset() // make it invalid / unusable
	t.toStringTabbed(n.left, tabLevel+1, printer)
	printer("\n")
	t.toStringTabbed(n.right, tabLevel+1, printer)
	printer("\n")
}

func (n *AVLTNode[K, V]) toStringTabbed(printer func(string)) {
	var sb strings.Builder
	sb.WriteString("Node [ height= ")
	sb.WriteString(strconv.Itoa(n.height))
	sb.WriteString(" ; key= <<")
	sb.WriteString(fmt.Sprint(n.keyVal.key))
	sb.WriteString(">> ;; size left= ")
	sb.WriteString(strconv.Itoa(n.sizeLeft))
	sb.WriteString(" ; size right= ")
	sb.WriteString(strconv.Itoa(n.sizeRight))
	sb.WriteString(" ;; father's key= <<")
	sb.WriteString(fmt.Sprint(n.father.keyVal.key))
	sb.WriteString(">> ... am I NIL? ")
	sb.WriteString(fmt.Sprint(n.father == n))
	sb.WriteString(" ]")
	printer(sb.String())
}

//

// PUBLIC functions

//

func NewAVLTree[K any, V any](keyZeroValue K, valueZeroValue V, keyExtractor KeyExtractor[K, V], comparator Comparator[K]) (*AVLTree[K, V], error) {
	if keyExtractor == nil {
		return nil, errors.New("key extractor must not be null")
	}
	if comparator == nil {
		return nil, errors.New("comparator must not be null")
	}
	t := new(AVLTree[K, V])
	t.keyZeroValue = keyZeroValue
	t.valueZeroValue = valueZeroValue
	t.comparator = comparator
	t.keyExtractor = keyExtractor
	t.size = 0
	t._NIL = nil
	_nil := t.newNode(keyZeroValue, valueZeroValue)
	t._NIL = _nil
	_nil.father = _nil
	_nil.left = _nil
	_nil.right = _nil
	_nil.prevInserted = _nil
	_nil.nextInserted = _nil
	_nil.height = DEPTH_INITIAL
	t.root = _nil

	// other optimization-based setup
	t.firstInserted = _nil
	t.minValue = _nil

	return t, nil
}

func (t *AVLTree[K, V]) Size() int {
	return t.size
}

func (t *AVLTree[K, V]) NILL() interface{} { return t._NIL }

func (t *AVLTree[K, V]) Put(key K, value V) (V, error) {
	n := t.newNode(key, value)
	return t.put(n)
}

func (t *AVLTree[K, V]) StringInto(printer func(string)) {
	printer("AVL Tree ")
	if t == nil {
		printer("\t- NULL!")
		return
	}
	if t.root == t._NIL {
		printer("\t- empty")
	}
	printer("of size= ")
	printer(strconv.Itoa(t.size))
	printer("; :\n")
	t.toStringTabbed(t.root, 0, printer)
}
func (t *AVLTree[K, V]) String() string {
	var sb strings.Builder
	t.StringInto(func(s string) { sb.WriteString(s) })
	return sb.String()
}

func (n *AVLTNode[K, V]) String() string {
	var sb strings.Builder
	n.toStringTabbed(func(s string) { sb.WriteString(s) })
	return sb.String()
}

//

func (p *KeyVal[K, V]) PairKey() K {
	return p.key
}
func (p *KeyVal[K, V]) PairValue() V {
	return p.value
}

// func (n *AVLTNode[K, V]) GetKeyVal() KeyVal[K, V] { return n.keyVal }

/**
 * Use with care.
 * <p>
 * {@inheritDoc}
 */
/****
 *
 *
 *
 *
@SuppressWarnings("unchecked")
@Override
protected V delete(nnn AVLTNode[K, V]) {
	boolean hasLeft, hasRight
	V v
	NodeAVL_Full nToBeDeleted, succMaybeDeleted
	if (root == NIL || nnn == NIL)
		return null
	v = null
	nToBeDeleted = (NodeAVL_Full) nnn
	v = nToBeDeleted.v
	if (size == 1 && comp.compare(root.k, nToBeDeleted.k) == 0) {
		v = super.delete(nToBeDeleted)
		firstInserted = (NodeAVL_Full) NIL
		((NodeAVL_Full) NIL).prevInserted = ((NodeAVL_Full) NIL).nextInserted = (NodeAVL_Full) NIL
		return v
	}
	// real deletion starts here:
	hasLeft = nToBeDeleted.left != NIL
	hasRight = nToBeDeleted.right != NIL
	succMaybeDeleted = hasRight ? (MapTreeAVLFull<K, V>.NodeAVL_Full) successorSorted(nnn) : //
			(MapTreeAVLFull<K, V>.NodeAVL_Full) (hasLeft ? predecessorSorted(nnn) : NIL)//

	v = super.delete(nnn)
	// adjust connections
	if (hasLeft || hasRight) {
		if (size == 1) {
			firstInserted = nToBeDeleted
			nToBeDeleted.nextInserted = nToBeDeleted.prevInserted = nToBeDeleted
		} else {
			// nnn wasn't the removed node ...
			// 1) unlink myself (nnn: nToBeDeleted) because that's me that should be removed
			// 2) then I re-link myself because I took the data held by the
			// node that has been removed in the end (succMaybeDeleted)
			// ..1) unlink myself
			nToBeDeleted.nextInserted.prevInserted = nToBeDeleted.prevInserted
			nToBeDeleted.prevInserted.nextInserted = nToBeDeleted.nextInserted
			// 2) then adjust my links to the really-removed-nodes ..
			nToBeDeleted.nextInserted = succMaybeDeleted.nextInserted
			nToBeDeleted.prevInserted = succMaybeDeleted.prevInserted
			// .. and the adjacent's nodes to point towards me
			nToBeDeleted.nextInserted.prevInserted = nToBeDeleted
			nToBeDeleted.prevInserted.nextInserted = nToBeDeleted
			if (succMaybeDeleted == firstInserted) { firstInserted = succMaybeDeleted.nextInserted }
		}
	} else {
		if (size == 1) {
			firstInserted = nToBeDeleted
			nToBeDeleted.nextInserted = nToBeDeleted.prevInserted = nToBeDeleted
		} else {
			if (nToBeDeleted == firstInserted)
				firstInserted = nToBeDeleted.nextInserted
			nToBeDeleted.nextInserted.prevInserted = nToBeDeleted.prevInserted
			nToBeDeleted.prevInserted.nextInserted = nToBeDeleted.nextInserted
		}
	}

	((NodeAVL_Full) NIL).nextInserted = ((NodeAVL_Full) NIL).prevInserted = (NodeAVL_Full) NIL
	if (root == NIL) {
		firstInserted = (NodeAVL_Full) NIL
		return v
	}
	return v
}



*/
