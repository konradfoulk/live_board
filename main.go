package main

import "fmt"

func main() {
	var myString = []rune("résumé")
	var indexed = myString[1]
	fmt.Printf("%v, %T\n", indexed, indexed)

	for i, v := range myString {
		fmt.Println(i, v)
	}

	fmt.Printf("\nThe length of 'myString' is %v\n", len(myString))

	var myRune = 'a'
	fmt.Printf("myRune = %v", myRune)

	var strSlice = []string{"s", "u", "b", "s", "c", "r", "i", "b", "e"}
	var catStr = ""
	for i := range strSlice {
		catStr += strSlice[i]
	}
	fmt.Printf("\n%v", catStr)
}
