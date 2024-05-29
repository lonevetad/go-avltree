package avltree

import (
	"fmt"
	"sync"
)

type ErrorAVLTree struct {
	errText string
}

func (et *ErrorAVLTree) Error() string {
	if et == nil {
		return "WAIT, THIS ERROR IS FRICKING NULL!"
	}
	return et.errText
}

func (et *ErrorAVLTree) String() string {
	if et == nil {
		return "ErrorAVLTree is nil"
	}
	return et.errText
}

func (et *ErrorAVLTree) IsNil() bool {
	return et == nil
}

func (et *ErrorAVLTree) IsEmpty() bool {
	return et == nil || et.errText == ""
}
func newErrorFromText(s string) *ErrorAVLTree {
	e := new(ErrorAVLTree)
	(*e).errText = s
	return e
}
func ne(o **sync.Once, constPointer **ErrorAVLTree, s string) *ErrorAVLTree {
	if *o == nil {
		once := new(sync.Once)
		*o = once
	}
	(*o).Do(func() {
		e := newErrorFromText(s)
		*constPointer = e
	})
	return *constPointer
}

// const ( KEY_NOT_FOUND *ErrorAVLTree = ne("key not found") )

// ERROR FROM PUBLIC FUNCTIONS

var __once_KEY_EXTRACTOR_NIL *sync.Once = nil
var __KEY_EXTRACTOR_NIL *ErrorAVLTree = nil

func KEY_EXTRACTOR_NIL() *ErrorAVLTree {
	return ne(&__once_KEY_EXTRACTOR_NIL, &__KEY_EXTRACTOR_NIL, //
		"key extractor is nil")
}

var __once_COMPARATOR_NIL *sync.Once = nil
var __ERR_COMPARATOR_NIL *ErrorAVLTree = nil

func COMPARATOR_NIL() *ErrorAVLTree {
	return ne(&__once_COMPARATOR_NIL, &__ERR_COMPARATOR_NIL, //
		"key comparator is nil")
}

// TODO

var __once_KEY_NOT_FOUND *sync.Once = nil
var __ERR_KEY_NOT_FOUND *ErrorAVLTree = nil

func KEY_NOT_FOUND() *ErrorAVLTree {
	return ne(&__once_KEY_NOT_FOUND, &__ERR_KEY_NOT_FOUND, //
		"key not found")
}

var __once_EMPTY_TREE *sync.Once = nil
var __ERR_EMPTY_TREE *ErrorAVLTree = nil

func EMPTY_TREE() *ErrorAVLTree {
	return ne(&__once_EMPTY_TREE, &__ERR_EMPTY_TREE, //
		"key not found")
}

var __once_VALUE_RETURNED_NIL *sync.Once = nil
var __ERR_VALUE_RETURNED_NIL *ErrorAVLTree = nil

func VALUE_RETURNED_NIL() *ErrorAVLTree {
	return ne(&__once_VALUE_RETURNED_NIL, &__ERR_VALUE_RETURNED_NIL, //
		"value returned is nil")
}

var __once_VALUE_RETURNED_NOT_NIL *sync.Once = nil
var __ERR_VALUE_RETURNED_NOT_NIL *ErrorAVLTree = nil

func VALUE_RETURNED_NOT_NIL[V any](value V) *ErrorAVLTree {
	return ne(&__once_VALUE_RETURNED_NOT_NIL, &__ERR_VALUE_RETURNED_NOT_NIL, //
		fmt.Sprintf("value returned should be nil but is:\n%v", value))
}

var __once_UNMATCHED_KEYS *sync.Once = nil
var __ERR_UNMATCHED_KEYS *ErrorAVLTree = nil

func UNMATCHED_KEYS(expected, actual int) *ErrorAVLTree {
	return ne(&__once_UNMATCHED_KEYS, &__ERR_UNMATCHED_KEYS, //
		fmt.Sprintf("actual, returned key (%d) is different than the expected one (%d)", actual, expected))
}

var __once_UNMATCHED_VALUES *sync.Once = nil
var __ERR_UNMATCHED_VALUES *ErrorAVLTree = nil

func UNMATCHED_VALUES(expected, actual string) *ErrorAVLTree {
	return ne(&__once_UNMATCHED_VALUES, &__ERR_UNMATCHED_VALUES, //
		fmt.Sprintf("actual, returned value (%s) is different than the expected one (%s)", actual, expected))
}
