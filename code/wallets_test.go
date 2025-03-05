package code

import "testing"

func TestWalletSerialization(t *testing.T) {
	wallets, _ := NewWallets()

	wallets.SaveToFile() // 触发序列化

	// 反序列化验证
	loaded, _ := NewWallets()
	if err := loaded.LoadFromFile(); err != nil {
		t.Fatalf("反序列化失败: %v", err)
	}
}
