package main

import "fmt"

func main() {
	s := "Привет"

	v := []rune(s)

	var result []rune

	result = append(result, v[:2]...)
	result = append(result, []rune("ВВ")...)
	result = append(result, v[2:]...)

	fmt.Println(string(result))
}
