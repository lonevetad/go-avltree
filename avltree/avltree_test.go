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

//
// TESTS
//

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

	err = testEqualityPrimitive(true, tree.Size(), 0, "size should be 0")
	if err != nil {
		t.Error(err)
	}
	if tree._NIL == nil {
		t.Errorf("the tree's \"_NIL\" should NOT be nil\n")
	}
	if tree.root != tree._NIL {
		t.Errorf("the tree should NOT have a root AND should be \"_NIL\"\n")
	}
	err = testIsLeaf(tree, tree._NIL)
	if err != nil {
		t.Error(err)
	}
	err = testNIL(tree, true, tree.root, "root is not _NIL")
	if err != nil {
		t.Error(err)
	}
	err = testNIL(tree, true, tree.root.father, "father is not _NIL")
	if err != nil {
		t.Error(err)
	}
	err = testEqualityPrimitive(true, tree._NIL.height, DEPTH_INITIAL, fmt.Sprintf("NIL's height should be: %d", DEPTH_INITIAL))
	if err != nil {
		t.Error(err)
	}
	err = testEqualityPrimitive(true, tree._NIL.sizeLeft, 0, "NIL's sizeLeft should be 0")
	if err != nil {
		t.Error(err)
	}
	err = testEqualityPrimitive(true, tree._NIL.sizeRight, 0, "NIL's sizeRight should be 0")
	if err != nil {
		t.Error(err)
	}
	err = testNIL(tree, true, tree.minValue, "minValue is not _NIL")
	if err != nil {
		t.Error(err)
	}
	err = testNIL(tree, true, tree.firstInserted, "firstInserted is not _NIL")
	if err != nil {
		t.Error(err)
	}
	err = testNIL(tree, true, tree.root.nextInOrder, "nextInOrder is not _NIL")
	if err != nil {
		t.Error(err)
	}
	err = testNIL(tree, true, tree.root.prevInOrder, "prevInOrder is not _NIL")
	if err != nil {
		t.Error(err)
	}
	err = testNIL(tree, true, tree.root.nextInserted, "nextInserted is not _NIL")
	if err != nil {
		t.Error(err)
	}
	err = testNIL(tree, true, tree.root.prevInserted, "prevInserted is not _NIL")
	if err != nil {
		t.Error(err)
	}
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
	err = testEqualityObj(true, oldData, tree.avlTreeConstructorParams.ValueZeroValue, EqualTestData, //
		fmt.Sprintf("putting a value on empty tree should return the \"value's zero-value\", but we have: %v", oldData))
	if err != nil {
		t.Error(err)
	}
	if tree.root == nil {
		t.Errorf("the tree's root should NOT be nil\n")
	}
	err = testIsLeaf(tree, tree._NIL)
	if err != nil {
		t.Error(err)
	}
	err = testNIL(tree, false, tree.root, "root is _NIL; should not be NIL")
	if err != nil {
		t.Error(err)
	}

	//

	err = testEqualityPrimitive(true, tree.Size(), 1, "size should be 1")
	if err != nil {
		t.Error(err)
	}

	if tree._NIL == nil {
		t.Errorf("the tree's \"_NIL\" should NOT be nil\n")
	}

	// internal nodes disposition

	err = testNIL(tree, true, tree.root.father, "father is not _NIL")
	if err != nil {
		t.Error(err)
	}
	err = testEqualityPrimitive(true, tree.root.height, 0, "new node's height should be: 0")
	if err != nil {
		t.Error(err)
	}
	err = testEqualityPrimitive(true, tree.root.sizeLeft, 0, "new node's sizeLeft should be 0")
	if err != nil {
		t.Error(err)
	}
	err = testEqualityPrimitive(true, tree.root.sizeRight, 0, "new node's sizeRight should be 0")
	if err != nil {
		t.Error(err)
	}

	err = testEqualityObj(true, tree.minValue, tree.root, EqualData, "min value node should be equal to root, since it's the only node here")
	if err != nil {
		t.Error(err)
	}
	err = testEqualityObj(true, tree.firstInserted, tree.root, EqualData, "first inserted node should be equal to root, since it's the only node here")
	if err != nil {
		t.Error(err)
	}

	err = testEqualityObj(true, tree.root.nextInOrder, tree.root, EqualData, "nextInOrder should loop to itself, i.e. root, since it's the only node here")
	if err != nil {
		t.Error(err)
	}
	err = testEqualityObj(true, tree.root.prevInOrder, tree.root, EqualData, "prevInOrder should loop to itself, i.e. root, since it's the only node here")
	if err != nil {
		t.Error(err)
	}
	err = testEqualityObj(true, tree.root.nextInserted, tree.root, EqualData, "nextInserted should loop to itself, i.e. root, since it's the only node here")
	if err != nil {
		t.Error(err)
	}
	err = testEqualityObj(true, tree.root.prevInserted, tree.root, EqualData, "prevInserted should loop to itself, i.e. root, since it's the only node here")
	if err != nil {
		t.Error(err)
	}
}

func Test_AddOne_2WithSameKey_Replace(t *testing.T) {

	tree, err := NewTree()
	if err != nil {
		t.Error(err)
	}

	data := NewTestData()
	data.Id = 0
	data.Text = "First"
	_, err = tree.Put(data.Id, data)
	if err != nil {
		t.Error(err)
	}

	d2 := NewTestData()
	d2.Id = data.Id
	d2.Text = "Second"
	oldData, err := tree.Put(d2.Id, d2)
	if err != nil {
		t.Error(err)
	}
	err = testEqualityObj(true, oldData, data, EqualTestData, "should be the first data-value inserted")
	if err != nil {
		t.Error(err)
	}
	err = testEqualityObj(true, tree.root.keyVal.value, d2, EqualTestData, "should be the second data-value inserted, since it should be replaced")
	if err != nil {
		t.Error(err)
	}

	//

	err = testEqualityPrimitive(true, tree.Size(), 1, "size should be 1")
	if err != nil {
		t.Error(err)
	}

	if tree._NIL == nil {
		t.Errorf("the tree's \"_NIL\" should NOT be nil\n")
	}

	// internal nodes disposition

	err = testNIL(tree, true, tree.root.father, "father is not _NIL")
	if err != nil {
		t.Error(err)
	}
	err = testEqualityPrimitive(true, tree.root.height, 0, "new node's height should be: 0")
	if err != nil {
		t.Error(err)
	}
	err = testEqualityPrimitive(true, tree.root.sizeLeft, 0, "new node's sizeLeft should be 0")
	if err != nil {
		t.Error(err)
	}
	err = testEqualityPrimitive(true, tree.root.sizeRight, 0, "new node's sizeRight should be 0")
	if err != nil {
		t.Error(err)
	}

	err = testEqualityObj(true, tree.minValue, tree.root, EqualData, "min value node should be equal to root, since it's the only node here")
	if err != nil {
		t.Error(err)
	}
	err = testEqualityObj(true, tree.firstInserted, tree.root, EqualData, "first inserted node should be equal to root, since it's the only node here")
	if err != nil {
		t.Error(err)
	}

	err = testEqualityObj(true, tree.root.nextInOrder, tree.root, EqualData, "nextInOrder should loop to itself, i.e. root, since it's the only node here")
	if err != nil {
		t.Error(err)
	}
	err = testEqualityObj(true, tree.root.prevInOrder, tree.root, EqualData, "prevInOrder should loop to itself, i.e. root, since it's the only node here")
	if err != nil {
		t.Error(err)
	}
	err = testEqualityObj(true, tree.root.nextInserted, tree.root, EqualData, "nextInserted should loop to itself, i.e. root, since it's the only node here")
	if err != nil {
		t.Error(err)
	}
	err = testEqualityObj(true, tree.root.prevInserted, tree.root, EqualData, "prevInserted should loop to itself, i.e. root, since it's the only node here")
	if err != nil {
		t.Error(err)
	}
}
func Test_AddOne_2WithSameKey_Ignore(t *testing.T) {

	tree, err := NewTree()
	tree.avlTreeConstructorParams.KeyCollisionBehavior = IgnoreInsertion
	if err != nil {
		t.Error(err)
	}

	data := NewTestData()
	data.Id = 0
	data.Text = "First"
	_, err = tree.Put(data.Id, data)
	if err != nil {
		t.Error(err)
	}

	d2 := NewTestData()
	d2.Id = data.Id
	d2.Text = "Second"
	oldData, err := tree.Put(d2.Id, d2)
	if err != nil {
		t.Error(err)
	}
	err = testEqualityObj(true, oldData, data, EqualTestData, "should be the first data-value inserted")
	if err != nil {
		t.Error(err)
	}
	err = testEqualityObj(true, tree.root.keyVal.value, data, EqualTestData, "should be the second data-value inserted, as if the Put would be rejected")
	if err != nil {
		t.Error(err)
	}

	//

	err = testEqualityPrimitive(true, tree.Size(), 1, "size should be 1")
	if err != nil {
		t.Error(err)
	}

	if tree._NIL == nil {
		t.Errorf("the tree's \"_NIL\" should NOT be nil\n")
	}

	// internal nodes disposition

	err = testNIL(tree, true, tree.root.father, "father is not _NIL")
	if err != nil {
		t.Error(err)
	}
	err = testEqualityPrimitive(true, tree.root.height, 0, "new node's height should be: 0")
	if err != nil {
		t.Error(err)
	}
	err = testEqualityPrimitive(true, tree.root.sizeLeft, 0, "new node's sizeLeft should be 0")
	if err != nil {
		t.Error(err)
	}
	err = testEqualityPrimitive(true, tree.root.sizeRight, 0, "new node's sizeRight should be 0")
	if err != nil {
		t.Error(err)
	}

	err = testEqualityObj(true, tree.minValue, tree.root, EqualData, "min value node should be equal to root, since it's the only node here")
	if err != nil {
		t.Error(err)
	}
	err = testEqualityObj(true, tree.firstInserted, tree.root, EqualData, "first inserted node should be equal to root, since it's the only node here")
	if err != nil {
		t.Error(err)
	}

	err = testEqualityObj(true, tree.root.nextInOrder, tree.root, EqualData, "nextInOrder should loop to itself, i.e. root, since it's the only node here")
	if err != nil {
		t.Error(err)
	}
	err = testEqualityObj(true, tree.root.prevInOrder, tree.root, EqualData, "prevInOrder should loop to itself, i.e. root, since it's the only node here")
	if err != nil {
		t.Error(err)
	}
	err = testEqualityObj(true, tree.root.nextInserted, tree.root, EqualData, "nextInserted should loop to itself, i.e. root, since it's the only node here")
	if err != nil {
		t.Error(err)
	}
	err = testEqualityObj(true, tree.root.prevInserted, tree.root, EqualData, "prevInserted should loop to itself, i.e. root, since it's the only node here")
	if err != nil {
		t.Error(err)
	}
}

// adding 2: [2,1], [2,3]

func Test_AddOne_2_InOrder(t *testing.T) {

	tree, err := NewTree()
	if err != nil {
		t.Error(err)
	}

	keys := []int{2, 1}
	var datas = make([]*TestData, len(keys))
	for i, k := range keys {

		data := NewTestData()
		data.Id = k
		data.Text = fmt.Sprintf("v_%d", i)
		datas[i] = data

		_, err = tree.Put(data.Id, data)
		if err != nil {
			t.Error(err)
		}
	}

	err = testNIL(tree, false, tree.root, "root should not _NIL")
	if err != nil {
		t.Error(err)
	}
	err = testEqualityObj(true, tree.root.keyVal.value, datas[0], EqualTestData, //
		fmt.Sprintf("root (%v) should be: %v", tree.root.keyVal.value, datas[0]))
	if err != nil {
		t.Error(err)
	}

	err = testNIL(tree, false, tree.root.left, "root's left should not _NIL")
	if err != nil {
		t.Error(err)
	}
	err = testEqualityObj(true, tree.root.left.keyVal.value, datas[1], EqualTestData, //
		fmt.Sprintf("root's left (%v) should be: %v", tree.root.left.keyVal.value, datas[1]))
	if err != nil {
		t.Error(err)
	}

	//

	expectSize := int64(len(datas))
	err = testEqualityPrimitive(true, tree.Size(), expectSize, fmt.Sprintf("size should be %d", expectSize))
	if err != nil {
		t.Error(err)
	}

	// internal nodes disposition

	err = testNIL(tree, true, tree.root.father, "father is not _NIL")
	if err != nil {
		t.Error(err)
	}
	err = testNIL(tree, true, tree.root.right, "root's right should be _NIL, but it's not")
	if err != nil {
		t.Error(err)
	}
	err = testEqualityPrimitive(true, tree.root.height, 1, "new node's height should be: 1")
	if err != nil {
		t.Error(err)
	}
	err = testEqualityPrimitive(true, tree.root.sizeLeft, 1, "new node's sizeLeft should be 1")
	if err != nil {
		t.Error(err)
	}
	err = testEqualityPrimitive(true, tree.root.sizeRight, 0, "new node's sizeRight should be 0")
	if err != nil {
		t.Error(err)
	}

	err = testEqualityObj(true, tree.minValue, tree.root.left, EqualData, "min value node should be equal to root's left")
	if err != nil {
		t.Error(err)
	}
	err = testEqualityObj(true, tree.firstInserted, tree.root, EqualData, "first inserted node should be equal to root")
	if err != nil {
		t.Error(err)
	}

	err = testEqualityObj(true, tree.root.nextInOrder, tree.root.left, EqualData, "root's nextInOrder should be root's left")
	if err != nil {
		t.Error(err)
	}
	err = testEqualityObj(true, tree.root.prevInOrder, tree.root.left, EqualData, "root's prevInOrder should be root's left")
	if err != nil {
		t.Error(err)
	}
	err = testEqualityObj(true, tree.root.nextInserted, tree.root.left, EqualData, "root's nextInserted should be root's left")
	if err != nil {
		t.Error(err)
	}
	err = testEqualityObj(true, tree.root.prevInserted, tree.root.left, EqualData, "root's prevInserted should be root's left")
	if err != nil {
		t.Error(err)
	}

	err = testEqualityObj(true, tree.root.left.father, tree.root, EqualData, "the second node's father should be root")
	if err != nil {
		t.Error(err)
	}
	err = testEqualityObj(true, tree.root.left.nextInOrder, tree.root, EqualData, "the second node's nextInOrder should be root")
	if err != nil {
		t.Error(err)
	}
	err = testEqualityObj(true, tree.root.left.prevInOrder, tree.root, EqualData, "the second node's prevInOrder should be root")
	if err != nil {
		t.Error(err)
	}
	err = testEqualityObj(true, tree.root.left.nextInserted, tree.root, EqualData, "the second node's nextInserted should be root")
	if err != nil {
		t.Error(err)
	}
	err = testEqualityObj(true, tree.root.left.prevInserted, tree.root, EqualData, "the second node's prevInserted should be root")
	if err != nil {
		t.Error(err)
	}
}

func Test_AddOne_2_ReverseOrder(t *testing.T) {

	tree, err := NewTree()
	if err != nil {
		t.Error(err)
	}

	keys := []int{2, 3}
	var datas = make([]*TestData, len(keys))
	for i, k := range keys {

		data := NewTestData()
		data.Id = k
		data.Text = fmt.Sprintf("v_%d", i)
		datas[i] = data

		_, err = tree.Put(data.Id, data)
		if err != nil {
			t.Error(err)
		}
	}

	err = testNIL(tree, false, tree.root, "root should not _NIL")
	if err != nil {
		t.Error(err)
	}
	err = testEqualityObj(true, tree.root.keyVal.value, datas[0], EqualTestData, //
		fmt.Sprintf("root (%v) should be: %v", tree.root.keyVal.value, datas[0]))
	if err != nil {
		t.Error(err)
	}

	err = testNIL(tree, false, tree.root.right, "root's right should not be _NIL")
	if err != nil {
		t.Error(err)
	}
	err = testEqualityObj(true, tree.root.right.keyVal.value, datas[1], EqualTestData, //
		fmt.Sprintf("root's right (%v) should be: %v", tree.root.right.keyVal.value, datas[1]))
	if err != nil {
		t.Error(err)
	}

	//

	expectSize := int64(len(datas))
	err = testEqualityPrimitive(true, tree.Size(), expectSize, fmt.Sprintf("size should be %d", expectSize))
	if err != nil {
		t.Error(err)
	}

	// internal nodes disposition

	err = testNIL(tree, true, tree.root.father, "father is not _NIL")
	if err != nil {
		t.Error(err)
	}
	err = testNIL(tree, true, tree.root.left, "root's left should be _NIL, but it's not")
	if err != nil {
		t.Error(err)
	}
	err = testEqualityPrimitive(true, tree.root.height, 1, "new node's height should be: 1")
	if err != nil {
		t.Error(err)
	}
	err = testEqualityPrimitive(true, tree.root.sizeLeft, 0, "new node's sizeleft should be 0")
	if err != nil {
		t.Error(err)
	}
	err = testEqualityPrimitive(true, tree.root.sizeRight, 1, "new node's sizeRight should be 1")
	if err != nil {
		t.Error(err)
	}

	err = testEqualityObj(true, tree.minValue, tree.root, EqualData, "min value node should be equal to root")
	if err != nil {
		t.Error(err)
	}
	err = testEqualityObj(true, tree.firstInserted, tree.root, EqualData, "first inserted node should be equal to root")
	if err != nil {
		t.Error(err)
	}

	err = testEqualityObj(true, tree.root.nextInOrder, tree.root.right, EqualData, "root's nextInOrder should be root's right")
	if err != nil {
		t.Error(err)
	}
	err = testEqualityObj(true, tree.root.prevInOrder, tree.root.right, EqualData, "root's prevInOrder should be root's right")
	if err != nil {
		t.Error(err)
	}
	err = testEqualityObj(true, tree.root.nextInserted, tree.root.right, EqualData, "root's nextInserted should be root's right")
	if err != nil {
		t.Error(err)
	}
	err = testEqualityObj(true, tree.root.prevInserted, tree.root.right, EqualData, "root's prevInserted should be root's right")
	if err != nil {
		t.Error(err)
	}

	err = testEqualityObj(true, tree.root.right.father, tree.root, EqualData, "the second node's father should be root")
	if err != nil {
		t.Error(err)
	}
	err = testEqualityObj(true, tree.root.right.nextInOrder, tree.root, EqualData, "the second node's nextInOrder should be root")
	if err != nil {
		t.Error(err)
	}
	err = testEqualityObj(true, tree.root.right.prevInOrder, tree.root, EqualData, "the second node's prevInOrder should be root")
	if err != nil {
		t.Error(err)
	}
	err = testEqualityObj(true, tree.root.right.nextInserted, tree.root, EqualData, "the second node's nextInserted should be root")
	if err != nil {
		t.Error(err)
	}
	err = testEqualityObj(true, tree.root.right.prevInserted, tree.root, EqualData, "the second node's prevInserted should be root")
	if err != nil {
		t.Error(err)
	}
}

//

type AddSetup struct {
	name string
	keys []int
	/* Predicted node order if the tree would be traversed in a "breadth first" way.
	As a clarifying example, the first value of this array, i.e. the number at index 0,
	refers to the tree's root. The stored values are the indexes of the "keys" above,
	which could be used to retrieve the aforementioned nodes.
	*/
	onBreadthVisit_IndexKeyData []int
	datas                       []*TestData
}

func Test_Add_3(t *testing.T) {

	/* adding 3: - 6 tests
	   -) no rotation in order: [2,1,3]
	   -) no rotation anti- order: [2,3,1]
	   -) "left left": [3,2,1]
	   - - albero prima della rotazione: ( ((1) -> 2 <- ()) -> 3 <- () )
	   - - albero dopo della rotazione: ((1) -> 2 <- (3))
	   -) "left right": [3,1,2] ... rotazione articolata
	   - - albero prima della rotazione: ( (() -> 1 <- (2)) -> 3 <- () )
	   - - albero dopo della rotazione: ((1) -> 2 <- (3))
	   -) "right right": [1,2,3]
	   - - albero prima della rotazione: ( () -> 1 <- (() -> 2 <- (3)) )
	   - - albero dopo della rotazione: ((1) -> 2 <- (3))
	   -) "right left": [1,3,2] ... rotazione articolata
	   - - albero prima della rotazione: ( () -> 1 <- ((2) -> 3 <- ()) )
	   - - albero dopo della rotazione: ((1) -> 2 <- (3))
	*/
	setups := []AddSetup{
		{"no rotation - in order", []int{2, 1, 3}, []int{
			0, 1, 2,
		}, nil},
		{"no rotation - reverse order", []int{2, 3, 1}, []int{
			0, 2, 1,
		}, nil},
		{"left left", []int{3, 2, 1}, []int{
			1, 2, 0, //2, 0, 1,
		}, nil},
		{"left right", []int{3, 1, 2}, []int{2, 1, 0}, nil},
		{"right right", []int{1, 2, 3}, []int{1, 0, 2}, nil},
		{"right left", []int{1, 3, 2}, []int{2, 0, 1}, nil},
	}

	for _, data := range setups {

		tree, err := NewTree()
		if err != nil {
			err = fmt.Errorf("on test {{\"%s\" - %v}} -- error: %s", data.name, data.keys, err)
			t.Error(err)
		}

		data.datas = make([]*TestData, len(data.keys))

		for i, id := range data.keys {
			dataTest := NewTestData()
			dataTest.Id = id
			dataTest.Text = fmt.Sprintf("v_%d", i)
			data.datas[i] = dataTest

			_, err = tree.Put(dataTest.Id, dataTest)
			if err != nil {
				err = fmt.Errorf("on test {{\"%s\" - %v}} -- error: %s", data.name, data.keys, err)
				t.Error(err)
			}
		}

		// early definitions
		var nodeNextInserted, nodePrevInserted *AVLTNode[int, *TestData]

		// root checks

		err = testNIL(tree, false, tree.root, "root should not _NIL")
		if err != nil {
			err = fmt.Errorf("on test {{\"%s\" - %v}} -- error: %s", data.name, data.keys, err)
			t.Error(err)
		}

		expectSize := int64(len(data.datas))
		err = testEqualityPrimitive(true, tree.Size(), expectSize, fmt.Sprintf("size should be %d", expectSize))
		if err != nil {
			err = fmt.Errorf("on test {{\"%s\" - %v}} -- error: %s", data.name, data.keys, err)
			t.Error(err)
		}

		indexRoot := 0
		rootsKeyIndex := data.onBreadthVisit_IndexKeyData[indexRoot]
		// indexLeft := data.onBreadthVisit_IndexKeyData[1]
		// indexRight := data.onBreadthVisit_IndexKeyData[2]
		dataRootExpected := data.datas[rootsKeyIndex]
		err = testEqualityObj(true, tree.root.keyVal.key, dataRootExpected.Id, EqualKey, //
			fmt.Sprintf("root key (%d) should be: %d", tree.root.keyVal.key, dataRootExpected.Id))
		if err != nil {
			err = fmt.Errorf("on test {{\"%s\" - %v}} -- error: %s", data.name, data.keys, err)
			t.Error(err)
		}
		err = testEqualityObj(true, tree.root.keyVal.value, dataRootExpected, EqualTestData, //
			fmt.Sprintf("root value (%v) should be: %v", tree.root.keyVal.value, dataRootExpected))
		if err != nil {
			err = fmt.Errorf("on test {{\"%s\" - %v}} -- error: %s", data.name, data.keys, err)
			t.Error(err)
		}

		err = testNIL(tree, true, tree.root.father, "father is not _NIL")
		if err != nil {
			err = fmt.Errorf("on test {{\"%s\" - %v}} -- error: %s", data.name, data.keys, err)
			t.Error(err)
		}

		err = testEqualityPrimitive(true, tree.root.height, 1, "new node's height should be: 1")
		if err != nil {
			err = fmt.Errorf("on test {{\"%s\" - %v}} -- error: %s", data.name, data.keys, err)
			t.Error(err)
		}
		err = testEqualityPrimitive(true, tree.root.sizeLeft, 1, "new node's sizeleft should be 1")
		if err != nil {
			err = fmt.Errorf("on test {{\"%s\" - %v}} -- error: %s", data.name, data.keys, err)
			t.Error(err)
		}
		err = testEqualityPrimitive(true, tree.root.sizeRight, 1, "new node's sizeRight should be 1")
		if err != nil {
			err = fmt.Errorf("on test {{\"%s\" - %v}} -- error: %s", data.name, data.keys, err)
			t.Error(err)
		}

		index_onBreadthVisit_IndexKeyData_root := 0
		index_onBreadthVisit_IndexKeyData_left := index_onBreadthVisit_IndexKeyData_root + 1
		index_onBreadthVisit_IndexKeyData_right := index_onBreadthVisit_IndexKeyData_root + 2

		indexKey_WhenRootWasAdded := data.onBreadthVisit_IndexKeyData[index_onBreadthVisit_IndexKeyData_root]
		dataRoot := data.datas[data.onBreadthVisit_IndexKeyData[index_onBreadthVisit_IndexKeyData_root]]
		rootRecalcolated := tree.getNode(dataRoot.Id)
		if rootRecalcolated != tree.root {
			t.Error("the tests are wrong! root has been wrongly indexed")
		}
		firstNodeInserted := tree.getNode(data.keys[0])

		indexRootNextInserted := data.onBreadthVisit_IndexKeyData[(indexKey_WhenRootWasAdded+1)%len(data.keys)]
		dataRootNextInserted := data.datas[indexRootNextInserted]
		tempIndexPrev := (indexKey_WhenRootWasAdded - 1)
		if tempIndexPrev < 0 {
			tempIndexPrev += len(data.keys)
		}
		indexRootPrevInserted := data.onBreadthVisit_IndexKeyData[tempIndexPrev]
		dataRootPrevInserted := data.datas[indexRootPrevInserted]

		nodeNextInserted = tree.getNode(dataRootNextInserted.Id)
		nodePrevInserted = tree.getNode(dataRootPrevInserted.Id)
		// secondNodeInserted = tree.getNode(data.keys[1])
		// thirdNodeInserted = tree.getNode(data.keys[2]) // 3 elements -> "3-1" then preceeds "0"

		/*
			firstNodeInserted := tree.getNode(dataRootExpected.Id)
			secondNodeNextInserted = tree.getNode(data.datas[indexLeft].Id)
			thirdNodePrevInserted = tree.getNode( data.datas[indexRight].Id ) // 3 elements -> "3-1" then preceeds "0"
		*/

		err = testEqualityObj(true, tree.minValue, tree.root.left, EqualData, "min value node should be equal to root 's left")
		if err != nil {
			err = fmt.Errorf("on test {{\"%s\" - %v}} -- error: %s", data.name, data.keys, err)
			t.Error(err)
		}
		err = testEqualityObj(true, tree.firstInserted, firstNodeInserted, EqualData, fmt.Sprintf("first inserted node (%v) should be equal to : %v", tree.firstInserted.keyVal.value, firstNodeInserted.keyVal.value))
		if err != nil {
			err = fmt.Errorf("on test {{\"%s\" - %v}} -- error: %s", data.name, data.keys, err)
			t.Error(err)
		}

		err = testEqualityObj(true, tree.root.nextInOrder, tree.root.right, EqualData, fmt.Sprintf("root's nextInOrder (whose value is: %v) should be root's right, with value: %v", tree.root.nextInOrder.keyVal.value, tree.root.right.keyVal.value))
		if err != nil {
			err = fmt.Errorf("on test {{\"%s\" - %v}} -- error: %s", data.name, data.keys, err)
			t.Error(err)
		}
		err = testEqualityObj(true, tree.root.prevInOrder, tree.root.left, EqualData, fmt.Sprintf("root's prevInOrder (whose value is: %v) should be root's left, with value: %v", tree.root.prevInOrder.keyVal.value, tree.root.left.keyVal.value))
		if err != nil {
			err = fmt.Errorf("on test {{\"%s\" - %v}} -- error: %s", data.name, data.keys, err)
			t.Error(err)
		}
		if err != nil {
			err = testEqualityObj(true, tree.root.nextInserted, nodeNextInserted, EqualData, fmt.Sprintf("root's nextInserted (whose value is: %v) should be the node with value: %v", tree.root.nextInserted.keyVal.value, nodeNextInserted.keyVal.value))
			err = fmt.Errorf("on test {{\"%s\" - %v}} -- error: %s", data.name, data.keys, err)
			t.Error(err)
		}
		err = testEqualityObj(true, tree.root.prevInserted, nodePrevInserted, EqualData,
			fmt.Sprintf("root's prevInserted (whose value is: %v) should be the node with value: %v (fetched with key: %d; index indexRootPrevInserted: %d, tempIndexPrev: %d, indexKey_WhenRootWasAdded:%d)",
				tree.root.prevInserted.keyVal.value, nodePrevInserted.keyVal.value,
				dataRootPrevInserted.Id, indexRootPrevInserted, tempIndexPrev, indexKey_WhenRootWasAdded))
		if err != nil {
			err = fmt.Errorf("on test {{\"%s\" - %v}} -- error: %s", data.name, data.keys, err)
			t.Error(err)
		}

		// root's left checks
		left := tree.root.left
		indexKey_WhenLeftWasAdded := data.onBreadthVisit_IndexKeyData[index_onBreadthVisit_IndexKeyData_left]
		// dataLeft := data.datas[data.onBreadthVisit_IndexKeyData[index_onBreadthVisit_IndexKeyData_left]]
		dataLeftNextInserted := data.datas[data.onBreadthVisit_IndexKeyData[(indexKey_WhenLeftWasAdded+1)%len(data.keys)]]
		tempIndexPrev = (indexKey_WhenLeftWasAdded - 1)
		if tempIndexPrev < 0 {
			tempIndexPrev += len(data.keys)
		}
		dataLeftPrevInserted := data.datas[data.onBreadthVisit_IndexKeyData[tempIndexPrev]]

		nodeNextInserted = tree.getNode(dataLeftNextInserted.Id)
		nodePrevInserted = tree.getNode(dataLeftPrevInserted.Id)

		err = testNIL(tree, false, left, "root's left should not be _NIL")
		if err != nil {
			err = fmt.Errorf("on test {{\"%s\" - %v}} -- error: %s", data.name, data.keys, err)
			t.Error(err)
		}
		err = testEqualityObj(true, left.father, tree.root, EqualData, fmt.Sprintf( //
			"left's father (left value: %v) should be root (value: %v), but we have as father: %v", //
			left.keyVal.value, tree.root.keyVal.value, left.father.keyVal.value))
		if err != nil {
			err = fmt.Errorf("on test {{\"%s\" - %v}} -- error: %s", data.name, data.keys, err)
			t.Error(err)
		}

		err = testNIL(tree, true, left.left, "left's left should be _NIL")
		if err != nil {
			err = fmt.Errorf("on test {{\"%s\" - %v}} -- error: %s", data.name, data.keys, err)
			t.Error(err)
		}
		err = testNIL(tree, true, left.right, "left's right should be _NIL")
		if err != nil {
			err = fmt.Errorf("on test {{\"%s\" - %v}} -- error: %s", data.name, data.keys, err)
			t.Error(err)
		}

		err = testEqualityObj(true, left.nextInOrder, tree.root, EqualData, fmt.Sprintf("root's left's nextInOrder should be root, with value: %v", tree.root.keyVal.value))
		if err != nil {
			err = fmt.Errorf("on test {{\"%s\" - %v}} -- error: %s", data.name, data.keys, err)
			t.Error(err)
		}
		err = testEqualityObj(true, left.prevInOrder, tree.root.right, EqualData, fmt.Sprintf("root's left's prevInOrder should be root's right, with value: %v", tree.root.right.keyVal.value))
		if err != nil {
			err = fmt.Errorf("on test {{\"%s\" - %v}} -- error: %s", data.name, data.keys, err)
			t.Error(err)
		}
		if err != nil {
			err = testEqualityObj(true, left.nextInserted, nodeNextInserted, EqualData, fmt.Sprintf("left's nextInserted should be the node with value: %v", nodeNextInserted.keyVal.value))
			err = fmt.Errorf("on test {{\"%s\" - %v}} -- error: %s", data.name, data.keys, err)
			t.Error(err)
		}
		err = testEqualityObj(true, left.prevInserted, nodePrevInserted, EqualData, fmt.Sprintf("left's prevInserted should be the node with value: %v", nodePrevInserted.keyVal.value))
		if err != nil {
			err = fmt.Errorf("on test {{\"%s\" - %v}} -- error: %s", data.name, data.keys, err)
			t.Error(err)
		}

		// root's right checks
		right := tree.root.right
		indexKey_WhenRightWasAdded := data.onBreadthVisit_IndexKeyData[index_onBreadthVisit_IndexKeyData_right]
		// dataLeft := data.datas[data.onBreadthVisit_IndexKeyData[index_onBreadthVisit_IndexKeyData_right]]
		dataRightNextInserted := data.datas[data.onBreadthVisit_IndexKeyData[(indexKey_WhenRightWasAdded+1)%len(data.keys)]]
		tempIndexPrev = (indexKey_WhenRightWasAdded - 1)
		if tempIndexPrev < 0 {
			tempIndexPrev += len(data.keys)
		}
		dataRightPrevInserted := data.datas[data.onBreadthVisit_IndexKeyData[tempIndexPrev]]

		nodeNextInserted = tree.getNode(dataRightNextInserted.Id)
		nodePrevInserted = tree.getNode(dataRightPrevInserted.Id)

		err = testNIL(tree, false, right, "root's right should not be _NIL")
		if err != nil {
			err = fmt.Errorf("on test {{\"%s\" - %v}} -- error: %s", data.name, data.keys, err)
			t.Error(err)
		}
		err = testEqualityObj(true, right.father, tree.root, EqualData, fmt.Sprintf( //
			"right's father (right value: %v) should be root (value: %v), but we have as father: %v", //
			right.keyVal.value, tree.root.keyVal.value, right.father.keyVal.value))
		if err != nil {
			err = fmt.Errorf("on test {{\"%s\" - %v}} -- error: %s", data.name, data.keys, err)
			t.Error(err)
		}

		err = testNIL(tree, true, right.left, "right's left should be _NIL")
		if err != nil {
			err = fmt.Errorf("on test {{\"%s\" - %v}} -- error: %s", data.name, data.keys, err)
			t.Error(err)
		}
		err = testNIL(tree, true, right.right, "right's right should be _NIL")
		if err != nil {
			err = fmt.Errorf("on test {{\"%s\" - %v}} -- error: %s", data.name, data.keys, err)
			t.Error(err)
		}

		err = testEqualityObj(true, right.nextInOrder, tree.root.left, EqualData, fmt.Sprintf("root's right's nextInOrder should be root's left, with value: %v", tree.root.left.keyVal.value))
		if err != nil {
			err = fmt.Errorf("on test {{\"%s\" - %v}} -- error: %s", data.name, data.keys, err)
			t.Error(err)
		}
		err = testEqualityObj(true, right.prevInOrder, tree.root, EqualData, fmt.Sprintf("root's right's prevInOrder should be root, with value: %v", tree.root.keyVal.value))
		if err != nil {
			err = fmt.Errorf("on test {{\"%s\" - %v}} -- error: %s", data.name, data.keys, err)
			t.Error(err)
		}
		if err != nil {
			err = testEqualityObj(true, right.nextInserted, nodeNextInserted, EqualData, fmt.Sprintf("right's nextInserted should be the node with value: %v", nodeNextInserted.keyVal.value))
			err = fmt.Errorf("on test {{\"%s\" - %v}} -- error: %s", data.name, data.keys, err)
			t.Error(err)
		}
		err = testEqualityObj(true, right.prevInserted, nodePrevInserted, EqualData, fmt.Sprintf("right's prevInserted should be the node with value: %v", nodePrevInserted.keyVal.value))
		if err != nil {
			err = fmt.Errorf("on test {{\"%s\" - %v}} -- error: %s", data.name, data.keys, err)
			t.Error(err)
		}
	}

}

//
//
//

func EqualKey[K int | int32 | int64](d1 K, d2 K) bool {
	return d1 == d2
}
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

func testEqualityObj[V any](shouldBeEqual bool, actual V, expected V, equalityPredicate func(V, V) bool, additionalErrorText string) error {
	shouldOrShouldNot := ""
	if !shouldBeEqual {
		shouldOrShouldNot = " NOT"
	}
	if shouldBeEqual != equalityPredicate(actual, expected) {
		err := fmt.Errorf("actual value %v should%s be equal to expected value %v; %s\n", actual, shouldOrShouldNot, expected, additionalErrorText)
		if err != nil {
			return err
		}
	}
	return nil
}

func testEqualityPrimitive[V int | int64 | int32](shouldBeEqual bool, actual V, expected V, additionalErrorText string) error {
	shouldOrShouldNot := ""
	if !shouldBeEqual {
		shouldOrShouldNot = " NOT"
	}
	if shouldBeEqual != (actual == expected) {
		err := fmt.Errorf("actual value %d should%s be equal to expected value %d; %s\n", actual, shouldOrShouldNot, expected, additionalErrorText)
		if err != nil {
			return err
		}
	}
	return nil
}
func testLinkEquality[K any, V any](shouldBeEqual bool, n *AVLTNode[K, V], nodeLinkedTo *AVLTNode[K, V], additionalErrorText string) error {
	shouldOrShouldNot := ""
	if !shouldBeEqual {
		shouldOrShouldNot = " NOT"
	}
	if shouldBeEqual != (n == nodeLinkedTo) {
		err := fmt.Errorf("Node should%s be equal ; %s ;\n\tkey node = %v\n\tkey link = %v\n", shouldOrShouldNot, additionalErrorText, n.keyVal.key, nodeLinkedTo.keyVal.key)
		if err != nil {
			return err
		}
	}
	return nil
}
func testNIL[K any, V any](tree *AVLTree[K, V], shouldBeNil bool, n *AVLTNode[K, V], additionalErrorText string) error {
	return testLinkEquality(shouldBeNil, n, tree._NIL, fmt.Sprintf("should be nil; %s", additionalErrorText))
}
func testIsLeaf[K any, V any](tree *AVLTree[K, V], n *AVLTNode[K, V]) error {
	err := testNIL(tree, true, n.left, "therefore, has a left")
	if err != nil {
		return err
	}
	err = testNIL(tree, true, n.right, "therefore, has a right")
	return err
}

func getNodePath[K any, V any](t *AVLTree[K, V], path []bool) (*AVLTNode[K, V], error) {
	n := t.root
	l := len(path)
	for i := 0; i < l && n != t._NIL; i++ {
		if path[i] {
			n = n.left
		} else {
			n = n.right
		}
	}
	if n == t._NIL { // error
		return nil, fmt.Errorf("node not found with path: %v", path)
	}
	return n, nil
}

func gnp[K any, V any](t *AVLTree[K, V], path []int) (*AVLTNode[K, V], error) {
	p := make([]bool, len(path))
	for i, isLeft := range path {
		p[i] = isLeft == 0
	}
	return getNodePath(t, p)
}

func DumpTreesForErrorsPrinter[K any, V any](t1 *AVLTree[K, V], t2 *AVLTree[K, V], additionalPreText string, printer func(s string)) {
	printer(additionalPreText)
	t1.StringInto(true, printer)
	printer("\nt2:\n")
	t2.StringInto(true, printer)
}
func DumpTreesForErrorsBuilder[K any, V any](t1 *AVLTree[K, V], t2 *AVLTree[K, V], additionalPreText string, sb strings.Builder) {
	DumpTreesForErrorsPrinter(t1, t2, additionalPreText, func(s string) { sb.WriteString(s) })
}
func DumpTreesForErrors[K any, V any](t1 *AVLTree[K, V], t2 *AVLTree[K, V], additionalPreText string) string {
	var sb strings.Builder
	DumpTreesForErrorsBuilder(t1, t2, additionalPreText, sb)
	return sb.String()
}
func CheckTrees[K any, V any](t1 *AVLTree[K, V], t2 *AVLTree[K, V]) (bool, error) {
	if t1 == t2 {
		return true, fmt.Errorf("indentitarly equal trees")
	}

	if t1.size != t2.size {
		return false, fmt.Errorf("different sizes: %d and %d\n", t1.size, t2.size)
	}

	if t1.IsEmpty() && t2.IsEmpty() {
		return true, nil
	}
	if t1.IsEmpty() != t2.IsEmpty() {
		if t1.IsEmpty() {
			return false, fmt.Errorf("t1 is empty but t2 is not: t2 has %d nodes", t2.size)
		}
		return false, fmt.Errorf("t1 is not empty but t2 is: t1 it has %d nodes", t2.size)
	}

	if t1.root.height != t2.root.height {
		return false, fmt.Errorf(DumpTreesForErrors(t1, t2, //
			fmt.Sprintf("they have different heights: t1's %d, t2's  %d\nt1:\n", t1.root.height, t2.root.height)))
	}

	maxHeight := t1.root.height
	pathRun := make([]bool, maxHeight)

	equal, err := checkTreesEquality(t1, t2, t1.root, t2.root, pathRun, 0)
	if (!equal) || (err != nil) {
		return false, err
	}

	// look for nodes held by "in order / inserted-chronological" pointers BUT without fathers, children
	// (or a father who has them as child, which is granted by the tree equality)

	forEaches := []avltree.ForEachMode{
		avltree.InOrder,
		avltree.ReverseInOrder,
		avltree.Queue,
		avltree.Stack,
	}

	for ife, fe := range forEaches {
		errs1 := make([]error)
		errs2 := make([]error)
		nodes_count1 := 0
		nodes_count2 := 0

		accumulator =
			func(isOne bool) {
				io := isOne
				return func(node *AVLTNode[K, V]) {
					// TODO : check for
					nodesToCheck_name := []string{"father", "left", "right"}
					nodes := []*AVLTNode[K, V]{node.father, node.left, node.right}

					// for i_n, X := range nodes{ ...
					// -) if X-node is NIL ->
					// - -) if(io) { nodes_count1++ } else { nodes_count2++ }
					// ... }
					// TODO proseguire anche con gli altri due
				}
			}
		t1.ForEach(fe, accumulator(1 == 1))
		t2.ForEach(fe, accumulator(2 == 1))

		// TODO - mettere checks sui "count" e sparare (stampare) gli errori in un unico "fmt.Errorf" da, poi, restituire
	}
	return true, nil
}

func composeErrorOnCheckTree[K any, V any](t1 *AVLTree[K, V], t2 *AVLTree[K, V], n1 *AVLTNode[K, V], n2 *AVLTNode[K, V], pathRun []bool, depthCurrent int, additionalText string) error {
	var sb strings.Builder
	var branchText string
	if pathRun[depthCurrent] {
		branchText = "left"
	} else {
		branchText = "right"
	}
	sb.WriteString(fmt.Sprintf("\twhile exploring %s branch at depth %d (complete path: %v), an error occour:\n\t", branchText, depthCurrent, pathRun))
	sb.WriteString(additionalText)
	sb.WriteString("\n\tdumping trees")

	// TODO: aggiungere print di un errore con anche il path ... e per finire, invoca: DumpTreesForErrors()
	return fmt.Errorf(sb.String())
}

/*
path true == left, false == right
*/
func checkTreesEquality[K any, V any](t1 *AVLTree[K, V], t2 *AVLTree[K, V], n1 *AVLTNode[K, V], n2 *AVLTNode[K, V], pathRun []bool, depthCurrent int) (bool, error) {
	//TODO: DEVO RIMUOVERE IL CODICE COMMENTATO PER CIò CHE NON è STRETTAMENTE "nil", PERCHè VIENE AFFRONTATO NEL CICLO SOTTO

	//  checking : if the nodes are strictly "nil" (NOT "t._NIL" !)

	if n1 == nil && n2 == nil || (n1 == t1._NIL && n2 == t2._NIL) {
		return true, nil
	}

	if n1 == nil /*|| n1 == t1._NIL*/ {
		// ERROR: SHOULD NOT BE NIL
		var nullity string = "null"
		/*if n1 == nil {
		nullity */ /*= "null"*/
		/*		} else {
				nullity = "t._NIL"
			}*/
		return false, composeErrorOnCheckTree(t1, t2, n1, n2, pathRun, depthCurrent, //
			fmt.Sprintf("node of first tree is %s (the second one didn't)", nullity))
	}
	if n2 == nil /* || n2 == t2._NIL */ {
		// ERROR: SHOULD NOT BE NIL
		var nullity string = "null"
		/*if n2 == nil {
		nullity*/ /*= "null"*/
		/*} else {
			nullity = "t._NIL"
		}*/
		return false, composeErrorOnCheckTree(t1, t2, n1, n2, pathRun, depthCurrent, //
			fmt.Sprintf("node of second tree is %s (the first one didn't)", nullity))
	}
	/*
		if (n1.father == t1._NIL) != (n2.father == t2._NIL) {
			return false, composeErrorOnCheckTree(t1, t2, n1, n2, pathRun, depthCurrent, //
				"the nodes have different father existance: one of them is nil, the other one does exists")
		}
	*/

	//  checking : height, size left, size right
	integers := []string{"height", "size left", "size right"}
	int_1 := []int64{n1.height, n1.sizeLeft, n1.sizeRight}
	int_2 := []int64{n2.height, n2.sizeLeft, n2.sizeRight}
	// TODO: fare i check e sparare gli errori

	var comp1, comp2 int

	// comparing node's keys

	/*
		keep track of the "pointers to pointers" in order to generalize the handling of
		the 4 nodes: the current node itself, its father, its left and its right.
		All of them are "pointers to nodes" whose node's key should be equal amont "1" and "2"
	*/
	pointersNode1 := []**AVLTNode[K, V]{&n1, &(n1.father), &(n1.left), &(n1.right)}
	pointersNode2 := []**AVLTNode[K, V]{&n2, &(n2.father), &(n2.left), &(n2.right)}
	pointerName := []string{"", "'s father", "'s left", "'s right"}

	// TODO: iterate
	i := 0
	l := len(pointerName)
	for ; i < l; i++ {
		pointer1 := pointersNode1[i]
		pointer2 := pointersNode2[i]
		nameNode := pointerName[i]

		node1 := *pointer1 // this node could be any among "pointerName"
		node2 := *pointer2

		// checking : NIL-ity

		if (node1 == t1._NIL) != (node2 == t2._NIL) { // the "XOR" ("^") does not exists, "!=" is equivalent
			return false, composeErrorOnCheckTree(t1, t2, n1, n2, pathRun, depthCurrent, //
				fmt.Sprintf("while comparing nodes%s NIL-ity, they are different: the nil-comparison results in < %t > for 1 and in < %t > for 2\n\t the checked node 1: %v\n\t the checked node 2: %v\n", //
					nameNode, (node1 == t1._NIL), (node2 == t2._NIL), node1, node2))
		}
		// if both nodes are NOT "NIL", then they will be checked in the for loop below

	}

	//  checking : keys, various
	// NOTE: those checks are shifted outside the for loop above because:
	// -) the father should already be checked (ERROR: IT WON'T IF N1 & N2 ARE THE ROOTS! HOW TO DEAL WITH THIS SITUATION [== "those roots"] ?)
	// -) the children (left & right) will be checked by the recursion -> no need to check tem __twice__ ["thrice", actually, due to the "father" check]
	// - -) "checking children" and "at first checking the nodes and THEN shifting the perspective onto the childer" are pratically the same thing:
	// .....the latter is just NOT reduntant
	whoseKeys := []string{"'s self", "nextInOrder", "prevInOrder", "nextInserted", "prevInserted"}
	keyOwners1 := []*AVLTNode[K, V]{n1, n1.nextInOrder, n1.prevInOrder, n1.nextInserted, n1.prevInserted}
	keyOwners2 := []*AVLTNode[K, V]{n2, n2.nextInOrder, n2.prevInOrder, n2.nextInserted, n2.prevInserted}
	i = 0
	l = len(whoseKeys)
	for ; i < l; i++ {
		keyOwnername := whoseKeys[i]
		node1 := keyOwners1[i]
		node2 := keyOwners2[i]

		comp1 = int(t1.avlTreeConstructorParams.Comparator(node1.keyVal.key, node2.keyVal.key))
		comp2 = int(t2.avlTreeConstructorParams.Comparator(node1.keyVal.key, node2.keyVal.key))

		if comp1 != 0 {
			return false, composeErrorOnCheckTree(t1, t2, n1, n2, pathRun, depthCurrent, //
				fmt.Sprintf("while comparing nodes%s key with tree 1 comparator, the comparison should be 0, but is: %d", keyOwnername, comp1))
		}
		if comp2 != 0 {
			return false, composeErrorOnCheckTree(t1, t2, n1, n2, pathRun, depthCurrent, //
				fmt.Sprintf("while comparing nodes%s key with tree 2 comparator, the comparison should be 0, but is: %d", keyOwnername, comp2))
		}
	}

	// ...... iterazione in corso, TODO: spostare tutta la roba sopra e sotto -> qui dentro
	// ...................... in teoria ho gia' fatto tutto ... eh eh

	// TODO con il 2

	// recursion on children

	pathRun[depthCurrent] = true
	equal, err := checkTreesEquality(t1, t2, n1.left, n1.left, pathRun, depthCurrent+1)
	if (!equal) || (err != nil) {
		// TODO : dump errors using path to show the "list" of nodes and choiches where this error occours
		return false, composeErrorOnCheckTree(t1, t2, n1, n2, pathRun, depthCurrent, "")
	}
	pathRun[depthCurrent] = false
	equal, err = checkTreesEquality(t1, t2, n1.right, n1.right, pathRun, depthCurrent+1)
	if (!equal) || (err != nil) {
		// TODO : dump errors using path to show the "list" of nodes and choiches where this error occours
		return false, composeErrorOnCheckTree(t1, t2, n1, n2, pathRun, depthCurrent, "")
	}

	return true, nil
}
