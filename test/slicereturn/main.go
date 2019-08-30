package main

import "fmt"

/*
返回一个slice切片


返回map中的value如果是切片的话，返回的是切片指针
*/

type(
	A struct {
		A string
	}
)

var (
	mm map[string][]*A
)

func main() {
	mm = make(map[string][]*A)
	mm["1"] = []*A{&A{"1"},&A{"2"}}


	tt := testReturnSlice()

	fmt.Printf("mm[\"1\"] = %p, tt = %p\n", mm["1"],tt)

	tt = []*A{&A{"bbbb"},&A{"232423"}}

	//如果想改变mm["1"]的值
	*tt[0] = A{"faffs"}
	fmt.Printf("mm[\"1\"] = %p, tt = %p\n", mm["1"],tt)
	fmt.Println(mm["1"][0])
	fmt.Println(tt[0])

}


func testReturnSlice()[]*A{
	return mm["1"]
}
