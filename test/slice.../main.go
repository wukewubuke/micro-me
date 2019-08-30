package main

import "fmt"

/*
函数传递切片和...的区别

*/

func main() {
	var aa  = []int{1,2,3,4,5}
	fmt.Printf("aa  ==>a的类型为:%T, aa的内存地址为:%p, aa指向的内存地址为:%p\n",aa, &aa, aa)
	tt(aa...)


	bb := []int{6,7,8}
	aa = append(aa,bb...)
	fmt.Println(aa)

}


func tt(t ...int){
	fmt.Printf("func tt ==>t的类型为:%T, t的内存地址为:%p, t指向的内存地址为:%p\n",t,&t,t)

	var i []int = t
	fmt.Printf("func i ==>t的类型为:%T, i的内存地址为:%p, i指向的内存地址为:%p\n",i,&i,i)
}

