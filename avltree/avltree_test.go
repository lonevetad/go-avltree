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

func NewMetadata(td *TestData) AVLTreeConstructorParams[int, *TestData] {
	avlTreeConstructorParams := AVLTreeConstructorParams[int, *TestData]{}
	avlTreeConstructorParams.KeyCollisionBehavior = Replace
	avlTreeConstructorParams.KeyZeroValue = td.Id
	avlTreeConstructorParams.ValueZeroValue = td
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

func TestNewTree(t *testing.T) {
	td := NewTestData()
	avlTreeConstructorParams := NewMetadata(td)

	tree, err := NewAVLTree(avlTreeConstructorParams)
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

	testInt(t, true, tree.Size(), 0, "")

	if tree._NIL == nil {
		t.Errorf("the tree's \"_NIL\" should NOT be nil\n")
	}
	if tree.root != tree._NIL {
		t.Errorf("the tree should NOT have a root AND should be \"_NIL\"\n")
	}
	testIsLeaf(t, tree, tree._NIL)
	testNIL(t, tree, true, tree.root, "root is not _NIL")
	testNIL(t, tree, true, tree.root.father, "father is not _NIL")
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

func testInt(t *testing.T, shouldBeEqual bool, actual int64, expected int64, additionalErrorText string) {
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
