package code

import (
	"bytes"
	"fmt"
	"testing"
)

func TestBase58(t *testing.T) {
	// 标准测试：输入包含前导零
	input := []byte{0x00, 0x00, 0x61}
	encoded := Base58Encode(input) // 预期 "11a"
	decoded := Base58Decode(encoded)
	fmt.Println(input)
	fmt.Println(decoded)
	if !bytes.Equal(decoded, input) {
		t.Errorf("前导零处理错误: 解码结果 %x", decoded)
	}

	// 非法字符测试
	invalidInput := []byte{'0'} // '0' 不在字母表中
	decoded = Base58Decode(invalidInput)
	if decoded != nil {
		t.Error("未拦截非法字符")
	}
}
