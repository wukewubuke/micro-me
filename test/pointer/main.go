package main

import "fmt"


/*
函数传递指针，实际也是值的传递
main函数中i的内存地址中保存的是a的地址
ttt函数中，aaa的内存地址和 main函数中i的内存地址是不通的
只不过他们保存的数据都是main函数中a的内存地址

*/

func main() {
	a:=100
	var i *int = &a
	ttt(i)



	fmt.Printf("a = %p, i = %p\n",&a, &i)
}




func ttt(aaa *int){
	fmt.Printf("ttt===> aaa = %p\n",&aaa)
}