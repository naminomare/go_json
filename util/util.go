package util

import "fmt"

// IfnilError , if err is not nil then panic(err)
func IfnilError(err error) {
	if err != nil {
		panic(err)
	}
}

// IfnilPrint , if err is not nil then fmt.Println(err)
func IfnilPrint(err error) {
	if err != nil {
		fmt.Println(err)
	}
}
