package charset

import (
	"fmt"
	"testing"
)

func TestDetermineEncoding(t *testing.T) {
	utf8 := []byte{230, 136, 145, 230, 152, 175, 85, 84, 70, 56}
	gbk := []byte{206, 210, 202, 199, 71, 66, 75}

	fmt.Println(UTF8To(GBK, string(utf8)))
	fmt.Println(ToUTF8(GBK, string(gbk)))

	//gbk := []byte{206, 210, 202, 199, 71, 66, 75}
	//a, name, certain := charset.DetermineEncoding(utf8, "")
	//fmt.Printf("编码：%v\n名称：%s\n确定：%t\n", a, name, certain)
}
