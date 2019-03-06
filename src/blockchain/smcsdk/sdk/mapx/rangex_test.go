package mapx

import (
	"fmt"
	"testing"
)

func TestForRange(t *testing.T) {
	m1 := make(map[int]string)
	m1[93] = "23231"
	m1[13] = "23423423234324"
	m1[54] = "3432432423"
	m1[23] = "3434545345345"
	ForRange(m1, func(k int, v string) {
		fmt.Printf("m1, key=%v value=%v\n", k, v)
	})

	fmt.Println("")

	m2 := make(map[string]int)
	m2["jshssfwerew"] = 23231
	m2["jjgjdldwer"] = 23423423234324
	m2["oeoiruwerw"] = 3432432423
	m2["iouoiudfs"] = 3434545345345
	ForRange(m2, func(k string, v int) {
		fmt.Printf("m2, key=%v value=%v\n", k, v)
	})

	fmt.Println("")

	m3 := make(map[bool]int)
	m3[true] = 23231
	m3[false] = 23423423234324
	m3[true] = 3432432423
	m3[false] = 3434545345345
	ForRange(m3, func(k bool, v int) {
		fmt.Printf("m3, key=%v value=%v\n", k, v)
	})

	fmt.Println("")

	m4 := make(map[float32]int)
	m4[1.2] = 23231
	m4[1.1314526] = 23423423234324
	m4[2.9] = 3432432423
	m4[0.9] = 3434545345345
	ForRange(m4, func(k float32, v int) {
		fmt.Printf("m3, key=%v value=%v\n", k, v)
	})
}
