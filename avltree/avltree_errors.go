package avltree

import "sync"

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

func ne(o **sync.Once, constPointer **ErrorAVLTree, s string) *ErrorAVLTree {
	if *o == nil {
		once := new(sync.Once)
		*o = once
	}
	(*o).Do(func() {
		e := new(ErrorAVLTree)
		(*e).errText = s
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
var __COMPARATOR_NIL *ErrorAVLTree = nil

func COMPARATOR_NIL() *ErrorAVLTree {
	return ne(&__once_COMPARATOR_NIL, &__COMPARATOR_NIL, //
		"key comparator is nil")
}

// TODO

var __once_KEY_NOT_FOUND *sync.Once = nil
var __KEY_NOT_FOUND *ErrorAVLTree = nil

func KEY_NOT_FOUND() *ErrorAVLTree {
	return ne(&__once_KEY_NOT_FOUND, &__KEY_NOT_FOUND, //
		"key not found")
}
