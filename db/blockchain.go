package db

import (
	"BTCTestProject/code"
	"BTCTestProject/pkg/pkg/mod/github.com/boltdb/bolt@v1.3.1"
	"fmt"
	"log"
)

const dbFile = "blockchain.db" //数据库文件名，当前目录下
const blockBucket = "blocks"   // 名称

type Blockchain struct {
	//Blocks []*code.Block // 一个数组，每个元素都是指针，存储block区块的地址
	Tip []byte
	db  *bolt.DB
}

type BlockchainIterator struct {
	currentHash []byte
	db          *bolt.DB
}

// 新建一个区块链
func NewBlockchain() *Blockchain {
	var tip []byte //存储区块链的二进制数据
	db, err := bolt.Open(dbFile, 0600, nil)
	if err != nil {
		log.Panic(err)
	}
	// 处理数据更新
	err = db.Update(func(tx *bolt.Tx) error {
		bucket := tx.Bucket([]byte(blockBucket)) //按照名称打开数据库表格
		if bucket == nil {
			fmt.Println("No blockchain, create a new blockchain")
			genesis := code.NewGenesisBlock()                  //创建创世区块
			bucket, err = tx.CreateBucket([]byte(blockBucket)) // 创建一个bucket
			if err != nil {
				log.Panic(err)
			}
			err = bucket.Put(genesis.Hash, genesis.Serialize()) // 存入数据
			if err != nil {
				log.Panic(err)
			}
			err = bucket.Put([]byte("1"), genesis.Hash) // 存入数据
			if err != nil {
				log.Panic(err)
			}
			tip = genesis.Hash
		} else {
			tip = bucket.Get([]byte("1"))
		}
		return nil
	})
	if err != nil {
		log.Panic(err)
	}
	bc := Blockchain{tip, db}
	return &bc
}

func (bc *Blockchain) AddBlock(data string) {
	var lastHash []byte
	err := bc.db.View(func(tx *bolt.Tx) error {
		b := tx.Bucket([]byte(blockBucket)) //取得数据
		lastHash = b.Get([]byte("1"))       //取得第一个块
		return nil
	})
	if err != nil {
		log.Panic(err)
	}
	newBlock := code.NewBlock(data, lastHash)
	err = bc.db.Update(func(tx *bolt.Tx) error {
		bucket := tx.Bucket([]byte(blockBucket))               //取出
		err := bucket.Put(newBlock.Hash, newBlock.Serialize()) //压入数据
		if err != nil {
			log.Panic(err)
		}
		err = bucket.Put([]byte("1"), newBlock.Hash)
		if err != nil {
			log.Panic(err)
		}
		bc.Tip = newBlock.Hash
		return nil
	})
}

func (bc *Blockchain) Iterator() *BlockchainIterator {
	bci := &BlockchainIterator{bc.Tip, bc.db}
	return bci //根据区块链创建区块链迭代器

}

func (it *BlockchainIterator) Next() *code.Block {
	var block *code.Block
	err := it.db.View(func(tx *bolt.Tx) error {
		b := tx.Bucket([]byte(blockBucket))
		encodedBlock := b.Get(it.currentHash) //获取二进制数据
		block, _ = code.Deserialize(encodedBlock)

		return nil
	})
	if err != nil {
		log.Panic(err)
	}
	it.currentHash = block.PrevBlockHash
	return block
}

//type Blockchain struct {
//	Blocks []*code.Block // 一个数组，每个元素都是指针，存储block区块的地址
//
//}

//// 增加一个区块
//func (blocks *Blockchain) AddBlock(data string) {
//	prevBlock := blocks.Blocks[len(blocks.Blocks)-1] // 取出最后一个区块
//	newBlock := code.NewBlock(data, prevBlock.Hash)  // 创建一个区块
//	blocks.Blocks = append(blocks.Blocks, newBlock)  //区块链插入新的区块
//}
//
//// 创建一个区块链
//func NewBlockchain() *Blockchain {
//	return &Blockchain{[]*code.Block{
//		code.NewGenesisBlock(),
//	}}
//}
