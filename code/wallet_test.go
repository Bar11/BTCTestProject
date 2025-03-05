package code

import (
	"testing"
)

func TestValidateAddress(t *testing.T) {
	validAddress := "1y43CQhpCFGv9zys5j55zcklibx9pct8M" // 中本聪地址
	if !ValidateAddress(validAddress) {
		t.Error("合法地址校验失败")
	}

	invalidAddress := "1A1zP1eP5QGefi2DMPTfTL5SLmv7DivfNx" // 篡改末位字符
	if ValidateAddress(invalidAddress) {
		t.Error("非法地址误判为有效")
	}
}
