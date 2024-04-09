package avltree_unit_test

import (
	"avltree_unit"
	"testing"

	avltree "github.com/lonevetad/go-avltree"
)

func TestNewTree(t *testing.T) {
	td := avltree_unit.NewTestData()
	avlTreeConstructorParams := avltree_unit.NewMetadata(td)

	tree, err := avltree.NewAVLTree(avlTreeConstructorParams)
	if err != nil {
		t.Error(err)
	}

	if tree == nil {
		t.Error("the new tree should not be nil")
	}
}
