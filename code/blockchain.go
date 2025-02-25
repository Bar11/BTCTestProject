package code

import (
	"encoding/hex"
	"fmt"
	"github.com/boltdb/bolt"
	"log"
	"os"
)

const dbFile = "blockchain.db" //数据库文件名，当前目录下
const blockBucket = "blocks"   // 名称
const genesisCoinbaseData = "genesis Coinbase Data"

type Blockchain struct {
	//Blocks []*code.Block // 一个数组，每个元素都是指针，存储block区块的地址
	Tip []byte
	DB  *bolt.DB
}

type BlockchainIterator struct {
	currentHash []byte
	db          *bolt.DB
}

// 新建一个区块链
func NewBlockchain(address string) *Blockchain {
	if dbExists() == false {
		fmt.Println("please create new db")
		os.Exit(1)
	}
	var tip []byte //存储区块链的二进制数据
	db, err := bolt.Open(dbFile, 0600, nil)
	if err != nil {
		log.Panic(err)
	}
	err = db.Update(func(tx *bolt.Tx) error {
		bucket := tx.Bucket([]byte(blockBucket))
		tip = bucket.Get([]byte("1"))
		return nil
	})

	bc := Blockchain{tip, db}
	return &bc
	//return CreateBlockChain(address)
}

// 创建一个区块链创建一个数据库
func CreateBlockChain(address string) *Blockchain {
	if dbExists() {
		fmt.Println("db already exists")
		os.Exit(1)
	}
	fmt.Println("creating new blockchain")
	var tip []byte //存储区块链的二进制数据
	db, err := bolt.Open(dbFile, 0600, nil)
	if err != nil {
		log.Panic(err)
	}
	err = db.Update(func(tx *bolt.Tx) error {
		cbtx := NewCoinBaseTX(address, genesisCoinbaseData)
		genesis := NewGenesisBlock(cbtx)
		b, err := tx.CreateBucket([]byte(blockBucket))
		if err != nil {
			log.Panic(err)
		}
		err = b.Put(genesis.Hash, genesis.Serialize())
		if err != nil {
			log.Panic(err)
		}
		err = b.Put([]byte("1"), genesis.Hash)
		if err != nil {
			log.Panic(err)
		}
		tip = genesis.Hash
		return nil
	})

	bc := Blockchain{tip, db}
	return &bc
}

// 挖矿带来的交易（）
func (bc *Blockchain) MineBlock(transactions []*Transaction) {
	var lastHash []byte //最后的哈希
	err := bc.DB.View(func(tx *bolt.Tx) error {
		b := tx.Bucket([]byte(blockBucket))
		lastHash = b.Get([]byte("1"))
		return nil
	})
	if err != nil {
		log.Panic(err)
	}
	newBlock := NewBlock(transactions, lastHash) // 创建一个新的区块
	err = bc.DB.Update(func(tx *bolt.Tx) error {
		b := tx.Bucket([]byte(blockBucket))
		err := b.Put(newBlock.Hash, newBlock.Serialize())
		if err != nil {
			log.Panic(err)
		}
		err = b.Put([]byte("1"), newBlock.Hash)
		if err != nil {
			log.Panic(err)
		}
		bc.Tip = newBlock.Hash //保存上一块的哈希
		return nil
	})
}

// 查找没有使用输出的交易列表
func (bc *Blockchain) FindUnSpendableOutPuts(address string) []Transaction {
	var unspendTXs []Transaction
	spentTXOS := make(map[string][]int) // 开辟内存
	bci := bc.Iterator()
	for {
		block := bci.Next()
		for _, tx := range block.Transactions {
			txID := hex.EncodeToString(tx.ID) // 获取交易ID
		Outputs:
			for outIndex, out := range tx.Vout {
				if spentTXOS[txID] != nil {
					for _, spentOut := range spentTXOS[txID] {
						if spentOut == outIndex {
							continue Outputs // 循环到不等位置
						}
					}
				}
				if out.CanBeUnlockedWith(address) {
					unspendTXs = append(unspendTXs, *tx) //加入列表
				}
			}
			if tx.IsCoinBase() == false {
				for _, in := range tx.Vin {
					if in.CanUnlockOutputWith(address) { //判断是否可以锁定
						inTxID := hex.EncodeToString(in.Txid)
						spentTXOS[inTxID] = append(spentTXOS[inTxID], in.Vout)
					}
				}
			}
		}
		if len(block.PrevBlockHash) == 0 { // 最后一块
			break
		}
	}
	return unspendTXs

}

// 获取所有没有使用的交易
func (bc *Blockchain) FindUTXO(address string) []TXOutput {
	var utxos []TXOutput
	unsentTransactions := bc.FindUnSpendableOutPuts(address)
	for _, tx := range unsentTransactions {
		for _, out := range tx.Vout {
			if out.CanBeUnlockedWith(address) { //是否锁定
				utxos = append(utxos, out)
			}
		}
	}
	return utxos
}

// 查找没有使用的输出以参考输入
func (bc *Blockchain) FindSpendableOutPuts(address string, amount int) (int, map[string][]int) {
	unspentoutputs := make(map[string][]int)         // 输出
	unspenttxs := bc.FindUnSpendableOutPuts(address) //根据地质查询所有交易
	accmulated := 0                                  // 累计
Work:
	for _, tx := range unspenttxs {
		txID := hex.EncodeToString(tx.ID)
		for outIndex, out := range tx.Vout {
			if out.CanBeUnlockedWith(address) && accmulated < amount {
				accmulated += out.Value
				unspentoutputs[txID] = append(unspentoutputs[txID], outIndex)
				if accmulated >= amount {
					break Work
				}
			}
		}
	}

	return accmulated, unspentoutputs
}

func (bc *Blockchain) AddBlock(transactions []*Transaction) {
	var lastHash []byte
	err := bc.DB.View(func(tx *bolt.Tx) error {
		b := tx.Bucket([]byte(blockBucket)) //取得数据
		lastHash = b.Get([]byte("1"))       //取得第一个块
		return nil
	})
	if err != nil {
		log.Panic(err)
	}
	newBlock := NewBlock(transactions, lastHash)
	err = bc.DB.Update(func(tx *bolt.Tx) error {
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
	bci := &BlockchainIterator{bc.Tip, bc.DB}
	return bci //根据区块链创建区块链迭代器

}

func (it *BlockchainIterator) Next() *Block {
	var block *Block
	err := it.db.View(func(tx *bolt.Tx) error {
		b := tx.Bucket([]byte(blockBucket))
		encodedBlock := b.Get(it.currentHash) //获取二进制数据
		block, _ = Deserialize(encodedBlock)

		return nil
	})
	if err != nil {
		log.Panic(err)
	}
	it.currentHash = block.PrevBlockHash
	return block
}

func dbExists() bool {
	if _, err := os.Stat(dbFile); os.IsNotExist(err) {
		return false
	}
	return true
}
