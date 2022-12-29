package blockchain

import (
	"database/sql"
	_ "github.com/mattn/go-sqlite3"
	"os"
	"time"
)

type Blockchain struct {
	Db    *sql.DB
	index uint64
}

func (chain *Blockchain) AddBlock(block *Block) error {
	chain.index += 1
	err := insertBlockToBlockchain(chain.Db, Base64Encode(block.CurrHash), serializeBlock(block))
	if err != nil {
		return err
	}

	return nil
}

func (chain *Blockchain) Size() uint64 {
	return getLastBlockId(chain.Db)
}

func (chain *Blockchain) Balance(address string) uint64 {
	var (
		balance uint64
		sblock  string
		block   *Block
	)
	rows, err := getPreviousBlocks(chain.Db, chain.index)
	if err != nil {
		return balance
	}
	defer rows.Close()
	for rows.Next() {
		rows.Scan(&sblock)
		block = deserializeBlock(sblock)
		if value, ok := block.Mapping[address]; ok {
			balance = value
			break
		}
	}

	return balance
}

func (chain *Blockchain) LastHash() []byte {
	return Base64Decode(getLastHash(chain.Db))
}

// Init Blockchain

const (
	genesisBlockHash = "genesis-block-hash"
	genesisReward    = 100
	storageName      = "storage-chain"
	storageValue     = 100
)

func NewChain(filename, receiver string) error {
	file, err := os.Create(filename)
	if err != nil {
		return err
	}
	file.Close()

	db, err := sql.Open("sqlite3", filename)
	if err != nil {
		return err
	}
	defer db.Close()

	err = createBlockchainTable(db)
	if err != nil {
		return err
	}

	chain := &Blockchain{
		Db: db,
	}

	genesis := &Block{
		CurrHash:  []byte(genesisBlockHash),
		Mapping:   make(map[string]uint64),
		Miner:     receiver,
		TimeStamp: time.Now().Format(time.RFC3339),
	}
	genesis.Mapping[storageName] = storageValue
	genesis.Mapping[receiver] = genesisReward

	err = chain.AddBlock(genesis)
	if err != nil {
		return err
	}

	return nil
}

// Load Blockchain

func LoadChain(filename string) *Blockchain {
	db, err := sql.Open("sqlite3", filename)
	if err != nil {
		return nil
	}

	chain := &Blockchain{
		Db: db,
	}
	chain.index = chain.Size()

	return chain
}
