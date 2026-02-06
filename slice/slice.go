package main

import "fmt"

func main() {
	a := []int{1, 2, 3, 4, 5}
	s := make([]int, 3, 5) // свой массив, len=3, cap=5
	copy(s, a[:3])         // s = [1 2 3], cap=5
	foo(s)                 // append пишет 9,9 в позиции [3] и [4] того же массива
	s = s[:5]              // расширяем len до 5, теперь видим [1 2 3 9 9]
	fmt.Println(a, s)
}

func foo(a []int) {
	a = append(a, 9, 9)
}
