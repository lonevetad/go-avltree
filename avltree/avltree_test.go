package avltree

import (
	"fmt"
	"strconv"
	"strings"
	"testing"
)

type TestData struct {
	Id   int
	Text string
}

func Extract(t *TestData) int {
	if t == nil {
		return -1
	}
	return t.Id
}

func IntCompare(i1 int, i2 int) int64 {
	if i1 > i2 {
		return 1
	}
	if i1 < i2 {
		return -1
	}
	return 0
}

func (td *TestData) String() string {
	if td == nil {
		return "null"
	}
	var sb strings.Builder
	sb.WriteString("<id= ")
	sb.WriteString(strconv.Itoa(td.Id))
	sb.WriteString("; text= \"")
	sb.WriteString(td.Text)
	sb.WriteString("\">")
	return sb.String()
}

func NewMetadata(keyZeroValue int, valueZeroValue *TestData) AVLTreeConstructorParams[int, *TestData] {
	avlTreeConstructorParams := AVLTreeConstructorParams[int, *TestData]{}
	avlTreeConstructorParams.KeyCollisionBehavior = Replace
	avlTreeConstructorParams.KeyZeroValue = keyZeroValue
	avlTreeConstructorParams.ValueZeroValue = valueZeroValue
	avlTreeConstructorParams.KeyExtractor = Extract
	avlTreeConstructorParams.Comparator = IntCompare
	return avlTreeConstructorParams
}

func NewTestData() *TestData {
	td := new(TestData)
	td.Id = -42
	td.Text = "HELLO NULL STRING"
	return td
}

func NewTree() (*AVLTree[int, *TestData], error) {
	//td := NewTestData()
	avlTreeConstructorParams := NewMetadata(-1000, nil)
	return NewAVLTree(avlTreeConstructorParams)
}

func TestNewTree(t *testing.T) {
	tree, err := NewTree()
	if err != nil {
		t.Error(err)
	}

	if tree == nil {
		t.Error("the new tree should not be nil\n")
		return
	}

	if tree.root == nil {
		t.Errorf("the tree's root should NOT be nil\n")
	}

	testEqualityPrimitive(t, true, tree.Size(), 0, "size should be 0")

	if tree._NIL == nil {
		t.Errorf("the tree's \"_NIL\" should NOT be nil\n")
	}
	if tree.root != tree._NIL {
		t.Errorf("the tree should NOT have a root AND should be \"_NIL\"\n")
	}
	testIsLeaf(t, tree, tree._NIL)
	testNIL(t, tree, true, tree.root, "root is not _NIL")
	testNIL(t, tree, true, tree.root.father, "father is not _NIL")
	testEqualityPrimitive(t, true, tree._NIL.height, DEPTH_INITIAL, fmt.Sprintf("NIL's height should be: %d", DEPTH_INITIAL))
	testEqualityPrimitive(t, true, tree._NIL.sizeLeft, 0, "NIL's sizeLeft should be 0")
	testEqualityPrimitive(t, true, tree._NIL.sizeRight, 0, "NIL's sizeRight should be 0")
	testNIL(t, tree, true, tree.minValue, "minValue is not _NIL")
	testNIL(t, tree, true, tree.firstInserted, "firstInserted is not _NIL")
	testNIL(t, tree, true, tree.root.nextInOrder, "nextInOrder is not _NIL")
	testNIL(t, tree, true, tree.root.prevInOrder, "prevInOrder is not _NIL")
	testNIL(t, tree, true, tree.root.nextInserted, "nextInserted is not _NIL")
	testNIL(t, tree, true, tree.root.prevInserted, "prevInserted is not _NIL")
	/*
		pektTree := tree.PeekInternalStructure_test()
		null := pektTree.Nil_Test()
		if pektTree.
	*/
}
func Test_AddOne(t *testing.T) {
	tree, err := NewTree()
	if err != nil {
		t.Error(err)
	}

	data := NewTestData()
	data.Id = 0
	data.Text = "First"
	oldData, err := tree.Put(data.Id, data)
	if err != nil {
		t.Error(err)
	}
	testEqualityObj(t, true, oldData, tree.avlTreeConstructorParams.ValueZeroValue, EqualTestData, //
		fmt.Sprintf("putting a value on empty tree should return the \"value's zero-value\", but we have: %v", oldData))

	if tree.root == nil {
		t.Errorf("the tree's root should NOT be nil\n")
	}
	testIsLeaf(t, tree, tree._NIL)
	testNIL(t, tree, false, tree.root, "root is _NIL; should not be NIL")

	//

	testEqualityPrimitive(t, true, tree.Size(), 1, "size should be 1")

	if tree._NIL == nil {
		t.Errorf("the tree's \"_NIL\" should NOT be nil\n")
	}

	// internal nodes disposition

	testNIL(t, tree, true, tree.root.father, "father is not _NIL")
	testEqualityPrimitive(t, true, tree.root.height, 0, "new node's height should be: 0")
	testEqualityPrimitive(t, true, tree.root.sizeLeft, 0, "new node's sizeLeft should be 0")
	testEqualityPrimitive(t, true, tree.root.sizeRight, 0, "new node's sizeRight should be 0")

	testEqualityObj(t, true, tree.minValue, tree.root, EqualData, "min value node should be equal to root, since it's the only node here")
	testEqualityObj(t, true, tree.firstInserted, tree.root, EqualData, "first inserted node should be equal to root, since it's the only node here")

	testEqualityObj(t, true, tree.root.nextInOrder, tree.root, EqualData, "nextInOrder should loop to itself, i.e. root, since it's the only node here")
	testEqualityObj(t, true, tree.root.prevInOrder, tree.root, EqualData, "prevInOrder should loop to itself, i.e. root, since it's the only node here")
	testEqualityObj(t, true, tree.root.nextInserted, tree.root, EqualData, "nextInserted should loop to itself, i.e. root, since it's the only node here")
	testEqualityObj(t, true, tree.root.prevInserted, tree.root, EqualData, "prevInserted should loop to itself, i.e. root, since it's the only node here")
}

//

func EqualData[V any](d1 *V, d2 *V) bool {
	return d1 == d2
}
func EqualTestData(d1 *TestData, d2 *TestData) bool {
	return EqualData[TestData](d1, d2)
}
func EqualTestDataDeep(d1 *TestData, d2 *TestData) bool {
	return d1 == d2 || ((d1 != nil) && (d2 != nil) && //
		(d1.Id == d2.Id) && (d1.Text == d2.Text))
}

func testEqualityObj[V any](t *testing.T, shouldBeEqual bool, actual V, expected V, equalityPredicate func(V, V) bool, additionalErrorText string) {
	shouldOrShouldNot := ""
	if !shouldBeEqual {
		shouldOrShouldNot = " NOT"
	}
	if shouldBeEqual != equalityPredicate(actual, expected) {
		t.Errorf("actual value %v should%s be equal to expected value %v; %s\n", actual, shouldOrShouldNot, expected, additionalErrorText)
	}
}

func testEqualityPrimitive[V int | int64 | int32](t *testing.T, shouldBeEqual bool, actual V, expected V, additionalErrorText string) {
	shouldOrShouldNot := ""
	if !shouldBeEqual {
		shouldOrShouldNot = " NOT"
	}
	if shouldBeEqual != (actual == expected) {
		t.Errorf("actual value %d should%s be equal to expected value %d; %s\n", actual, shouldOrShouldNot, expected, additionalErrorText)
	}
}
func testLinkEquality[K any, V any](t *testing.T, shouldBeEqual bool, n *AVLTNode[K, V], nodeLinkedTo *AVLTNode[K, V], additionalErrorText string) {
	shouldOrShouldNot := ""
	if !shouldBeEqual {
		shouldOrShouldNot = " NOT"
	}
	if shouldBeEqual != (n == nodeLinkedTo) {
		t.Errorf("Node should%s be equal ; %s ;\n\tkey node = %v\n\tkey link = %v\n", shouldOrShouldNot, additionalErrorText, n.keyVal.key, nodeLinkedTo.keyVal.key)
	}
}
func testNIL[K any, V any](t *testing.T, tree *AVLTree[K, V], shouldBeNil bool, n *AVLTNode[K, V], additionalErrorText string) {
	testLinkEquality(t, shouldBeNil, n, tree._NIL, fmt.Sprintf("should be nil; %s", additionalErrorText))
}
func testIsLeaf[K any, V any](t *testing.T, tree *AVLTree[K, V], n *AVLTNode[K, V]) {
	testNIL(t, tree, true, n.left, "therefore, has a left")
	testNIL(t, tree, true, n.right, "therefore, has a right")
}
