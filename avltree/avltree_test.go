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
	name            string
	keys            []int
	orderingBreadth []int
	datas           []*TestData
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
			2, 0, 1,
		}, nil},
		{"left right", []int{3, 1, 2}, []int{2, 1, 0}, nil},
		{"right right", []int{1, 2, 3}, []int{1, 0, 2}, nil},
		{"left left", []int{1, 3, 2}, []int{1, 2, 0}, nil},
	}

	for i_test, data := range setups {

		tree, err := NewTree()
		if err != nil {
			err = fmt.Errorf("on test \"%s\" -- error: %s", err)
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
				err = fmt.Errorf("on test \"%s\" -- error: %s", err)
				t.Error(err)
			}
		}

		// early definitions
		var nodeNextInserted, nodePrevInserted *AVLTNode[int, *TestData]

		// root checks

		err = testNIL(tree, false, tree.root, "root should not _NIL")
		if err != nil {
			err = fmt.Errorf("on test \"%s\" -- error: %s", err)
			t.Error(err)
		}
		indexRoot := data.orderingBreadth[0]
		indexLeft := data.orderingBreadth[1]
		indexRight := data.orderingBreadth[2]
		dataRootExpected := data.datas[indexRoot]
		err = testEqualityObj(true, tree.root.keyVal.key, dataRootExpected.Id, EqualKey, //
			fmt.Sprintf("root key (%d) should be: %d", tree.root.keyVal.key, dataRootExpected.Id, indexRoot))
		if err != nil {
			err = fmt.Errorf("on test \"%s\" -- error: %s", err)
			t.Error(err)
		}
		err = testEqualityObj(true, tree.root.keyVal.value, dataRootExpected, EqualTestData, //
			fmt.Sprintf("root value (%v) should be: %v", tree.root.keyVal.value, dataRootExpected))
		if err != nil {
			err = fmt.Errorf("on test \"%s\" -- error: %s", err)
			t.Error(err)
		}

		err = testNIL(tree, true, tree.root.father, "father is not _NIL")
		if err != nil {
			err = fmt.Errorf("on test \"%s\" -- error: %s", err)
			t.Error(err)
		}

		err = testEqualityPrimitive(true, tree.root.height, 1, "new node's height should be: 1")
		if err != nil {
			err = fmt.Errorf("on test \"%s\" -- error: %s", err)
			t.Error(err)
		}
		err = testEqualityPrimitive(true, tree.root.sizeLeft, 1, "new node's sizeleft should be 1")
		if err != nil {
			err = fmt.Errorf("on test \"%s\" -- error: %s", err)
			t.Error(err)
		}
		err = testEqualityPrimitive(true, tree.root.sizeRight, 1, "new node's sizeRight should be 1")
		if err != nil {
			err = fmt.Errorf("on test \"%s\" -- error: %s", err)
			t.Error(err)
		}

		firstNodeInserted := tree.getNode(data.keys[0])
		nodeNextInserted = tree.getNode(data.keys[1])
		nodePrevInserted = tree.getNode(data.keys[2]) // 3 elements -> "3-1" then preceeds "0"

		err = testEqualityObj(true, tree.minValue, tree.root.left, EqualData, "min value node should be equal to root 's left")
		if err != nil {
			err = fmt.Errorf("on test \"%s\" -- error: %s", err)
			t.Error(err)
		}
		err = testEqualityObj(true, tree.firstInserted, firstNodeInserted, EqualData, "first inserted node should be equal to root")
		if err != nil {
			err = fmt.Errorf("on test \"%s\" -- error: %s", err)
			t.Error(err)
		}

		err = testEqualityObj(true, tree.root.nextInOrder, nodeNextInserted, EqualData, fmt.Sprintf("root's nextInOrder should be the node with value: %v", nodeNextInserted.keyVal.value))
		if err != nil {
			err = fmt.Errorf("on test \"%s\" -- error: %s", err)
			t.Error(err)
		}

		// TODO PROSEGIURE DA QUIIIIIIIIIIIIIIIIIIIIIIIII I "cosa" PASSARE PER PARAMETRO (dei checks qua sotto)

		err = testEqualityObj(true, tree.root.prevInOrder, tree.root.left, EqualData, "root's prevInOrder should be root's left")
		if err != nil {
			err = fmt.Errorf("on test \"%s\" -- error: %s", err)
			t.Error(err)
		}
		err = testEqualityObj(true, tree.root.nextInserted, tree.root.left, EqualData, "root's nextInserted should be root's left")
		if err != nil {
			err = fmt.Errorf("on test \"%s\" -- error: %s", err)
			t.Error(err)
		}
		err = testEqualityObj(true, tree.root.prevInserted, tree.root.right, EqualData, "root's prevInserted should be root's right")
		if err != nil {
			err = fmt.Errorf("on test \"%s\" -- error: %s", err)
			t.Error(err)
		}

		// root's left checks

		err = testNIL(tree, false, tree.root.left, "root's left should not be _NIL")
		if err != nil {
			err = fmt.Errorf("on test \"%s\" -- error: %s", err)
			t.Error(err)
		}

		// TODO: spostare il codice sotto (da fuori dal ciclo a dentro il ciclo)
		// per impostare test e
	}

	err = testEqualityObj(true, tree.root.left.keyVal.value, datas[1], EqualTestData, //
		fmt.Sprintf("root's left (%v) should be: %v", tree.root.left.keyVal.value, datas[1]))
	if err != nil {
		err = fmt.Errorf("on test \"%s\" -- error: %s", err)
		t.Error(err)
	}
	err = testNIL(tree, false, tree.root.right, "root's right should not be _NIL")
	if err != nil {
		err = fmt.Errorf("on test \"%s\" -- error: %s", err)
		t.Error(err)
	}
	err = testEqualityObj(true, tree.root.right.keyVal.value, datas[2], EqualTestData, //
		fmt.Sprintf("root's right (%v) should be: %v", tree.root.right.keyVal.value, datas[2]))
	if err != nil {
		err = fmt.Errorf("on test \"%s\" -- error: %s", err)
		t.Error(err)
	}

	//

	expectSize := int64(len(datas))
	err = testEqualityPrimitive(true, tree.Size(), expectSize, fmt.Sprintf("size should be %d", expectSize))
	if err != nil {
		err = fmt.Errorf("on test \"%s\" -- error: %s", err)
		t.Error(err)
	}

	//

	//

	err = testEqualityObj(true, tree.root.left.father, tree.root, EqualData, "the left node's father should be root")
	if err != nil {
		err = fmt.Errorf("on test \"%s\" -- error: %s", err)
		t.Error(err)
	}
	err = testEqualityObj(true, tree.root.left.nextInOrder, tree.root, EqualData, "the left node's nextInOrder should be root")
	if err != nil {
		err = fmt.Errorf("on test \"%s\" -- error: %s", err)
		t.Error(err)
	}
	err = testEqualityObj(true, tree.root.left.prevInOrder, tree.root.right, EqualData, "the left node's prevInOrder should be root's right")
	if err != nil {
		err = fmt.Errorf("on test \"%s\" -- error: %s", err)
		t.Error(err)
	}
	err = testEqualityObj(true, tree.root.left.nextInserted, tree.root.right, EqualData, "the left node's nextInserted should be root")
	if err != nil {
		err = fmt.Errorf("on test \"%s\" -- error: %s", err)
		t.Error(err)
	}
	err = testEqualityObj(true, tree.root.left.prevInserted, tree.root, EqualData, "the left node's prevInserted should be root")
	if err != nil {
		err = fmt.Errorf("on test \"%s\" -- error: %s", err)
		t.Error(err)
	}

	err = testEqualityObj(true, tree.root.right.father, tree.root, EqualData, "the right node's father should be root")
	if err != nil {
		err = fmt.Errorf("on test \"%s\" -- error: %s", err)
		t.Error(err)
	}
	err = testEqualityObj(true, tree.root.right.nextInOrder, tree.root.left, EqualData, "the right node's nextInOrder should be root's left")
	if err != nil {
		err = fmt.Errorf("on test \"%s\" -- error: %s", err)
		t.Error(err)
	}
	err = testEqualityObj(true, tree.root.right.prevInOrder, tree.root, EqualData, "the right node's prevInOrder should be root")
	if err != nil {
		err = fmt.Errorf("on test \"%s\" -- error: %s", err)
		t.Error(err)
	}
	err = testEqualityObj(true, tree.root.right.nextInserted, tree.root, EqualData, "the right node's nextInserted should be root")
	if err != nil {
		err = fmt.Errorf("on test \"%s\" -- error: %s", err)
		t.Error(err)
	}
	err = testEqualityObj(true, tree.root.right.prevInserted, tree.root.left, EqualData, "the right node's prevInserted should be root's left")
	if err != nil {
		err = fmt.Errorf("on test \"%s\" -- error: %s", err)
		t.Error(err)
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
