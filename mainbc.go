package main

import (
	bc "blockchain/blockchain"
	"fmt"
)

const (
	dbName = "blockchain.db"
)

func main() {
	miner := bc.NewUser()
	err := bc.NewChain(dbName, miner.Address())
	if err != nil {
		fmt.Println(err.Error())
		return
	}

	chain := bc.LoadChain(dbName)
	defer chain.Db.Close()

	for i := 0; i < 3; i++ {
		block := bc.NewBlock(miner.Address(), chain.LastHash())
		block.AddTransaction(chain, bc.NewTransaction(miner, "aaa", 5, chain.LastHash()))
		block.AddTransaction(chain, bc.NewTransaction(miner, "bbbb", 3, chain.LastHash()))
		block.Accept(chain, miner, make(chan bool))
		chain.AddBlock(block)
	}

	var sblock string
	rows, err := bc.GetBlocks(chain.Db)
	if err != nil {
		panic(err)
	}
	for rows.Next() {
		rows.Scan(&sblock)
		fmt.Println(sblock)
	}
}
