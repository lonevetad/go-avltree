package avltree_test_manual

import (
	"fmt"
	"strconv"
	"strings"

	avltree "github.com/lonevetad/go-avltree"
)

type TestData struct {
	id   int
	text string
}

func extract(t *TestData) int {
	if t == nil {
		return -1
	}
	return t.id
}

func compare(i1 int, i2 int) int64 {
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
	sb.WriteString(strconv.Itoa(td.id))
	sb.WriteString("; text= \"")
	sb.WriteString(td.text)
	sb.WriteString("\">")
	return sb.String()
}

func printAVLTree(t *avltree.AVLTree[int, *TestData]) {
	fmt.Printf("t:\n%v\n;\n\n and t's NIL:\n %v\n", t, t.NILL())
}

func printForEach(id int, td *TestData) {
	fmt.Printf("%s, ", td.String())
}

//

func main() {
	td := new(TestData)
	td.id = -666
	td.text = "HELLO NULL STRING"
	fmt.Printf("null data: %v\n", td)
	avlTreeConstructorParams := avltree.AVLTreeConstructorParams[int, *TestData]{}
	avlTreeConstructorParams.KeyCollisionBehavior = avltree.Replace
	avlTreeConstructorParams.KeyZeroValue = td.id
	avlTreeConstructorParams.ValueZeroValue = td
	avlTreeConstructorParams.KeyExtractor = extract
	avlTreeConstructorParams.Comparator = compare

	forEaches := []avltree.ForEachMode{
		avltree.InOrder,
		avltree.ReverseInOrder,
		avltree.Queue,
		avltree.Stack,
	}
	datasets_k := [][]int{
		{88, 7, 100}, //
		{5, 4, 3, 2, 100, 1, 80, 1111, 44, 22, 99, 84, 83},
		{1, 2, 3, 4, 5, 6, 7, 8, 9, 0, 10, 11, 12, 13, 14, 15, 16},
		{2, 2, 3, 3, 40, 40, 1, 1, 33, 33, 55, 55, 37, 37},
	}
	datasets_v := [][]string{
		{"Adam", "Eevee", "GOD HIMSELF"},
		nil,
		nil,
		nil,
	}

	for i, ids := range datasets_k {
		texts := datasets_v[i]
		fmt.Printf("\n\n\n -------------------------------------------------------------------\n")
		fmt.Printf("beginning the cycle # %d\n", i)
		t, err := avltree.NewAVLTree(avlTreeConstructorParams)
		if err != nil {
			fmt.Print("ERROR!")
			fmt.Print(err)
			return
		}
		for i, id := range ids {
			td_temp := new(TestData)
			td_temp.id = id
			if texts == nil {
				td_temp.text = fmt.Sprintf("t_%d", i)
			} else {
				td_temp.text = texts[i]
			}

			fmt.Printf("\n\n putting the %d-ish item: %v\n", i, td_temp)
			t.Put(id, td_temp)
			printAVLTree(t)
			fmt.Println("-------\ntesting all for-eaches:")
			for ife, fe := range forEaches {
				fmt.Printf("- - for-eacher #%d : %d\n\t =", ife, fe)
				t.ForEach(fe, printForEach)
				fmt.Println()
			}
		}

		fmt.Println("\n now removing key: 3")
		oldVal, err := t.Remove(3)
		if err != nil {
			fmt.Printf("ERROR:  %s\n", err)
		} else {
			fmt.Printf("deleted value %v\n", oldVal)
			printAVLTree(t)
			fmt.Println()
			fmt.Println("-------\ntesting all for-eaches:")
			for ife, fe := range forEaches {
				fmt.Printf("- - for-eacher #%d : %d\n\t =", ife, fe)
				err = t.ForEach(fe, printForEach)
				if err != nil {
					fmt.Println(err)
				}
				fmt.Println()
			}
			fmt.Println()
		}
		fmt.Println()
	}

	fmt.Println("finish")
}
