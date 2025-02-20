package main

import (
	"BTCTestProject/cmd"
	"BTCTestProject/db"
)

//TIP <p>To run your code, right-click the code and select <b>Run</b>.</p> <p>Alternatively, click
// the <icon src="AllIcons.Actions.Execute"/> icon in the gutter and select the <b>Run</b> menu item from here.</p>

func main() {

	//fmt.Println("Hello World")
	//bc := db.NewBlockchain()
	//bc.AddBlock("zhangsan pay lisi 10")
	//bc.AddBlock("zhangsan pay lisi 20")
	//bc.AddBlock("zhangsan pay lisi 30")
	//
	//for _, block := range bc.Blocks {
	//	fmt.Printf("timestamp:%d\n", block.Timestamp)
	//	fmt.Printf("PrevBlockHash:%x\n", block.PrevBlockHash)
	//	fmt.Printf("data:%s\n", block.Data)
	//	fmt.Printf("Hash:%x\n", block.Hash)
	//	pow := code.NewProofOfWork(block)
	//	fmt.Printf("pow %s\n", strconv.FormatBool(pow.Validate()))
	//	fmt.Println()
	//}
	block := db.NewBlockchain() //创建区块链
	defer block.DB.Close()      // 关闭数据库
	cli := cmd.CLI{block}       //创建命令行
	cli.Run()

}
