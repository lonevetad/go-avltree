module avltree_unit_test

go 1.22.1

require github.com/lonevetad/go-avltree v0.0.0

require "avltree_unit" v0.0.0

replace github.com/lonevetad/go-avltree v0.0.0 => ../../avltree

replace "avltree_unit" v0.0.0 => ../avltree_unit
