package blockchain

import (
	"bytes"
	"crypto/rsa"
	"errors"
	"math/big"
	"sort"
	"time"
)

type Block struct {
	CurrHash     []byte
	PrevHash     []byte
	Nonce        uint64
	Difficulty   uint8
	Miner        string
	Signature    string
	TimeStamp    string
	Transactions []Transaction
	Mapping      map[string]uint64
}

const (
	transactionsCountLimit = 2
	storageChain           = "storage-chain"
)

func (block *Block) AddTransaction(chain *Blockchain, tx *Transaction) error {
	if tx == nil {
		return errors.New("tx is null")
	}
	if tx.Value == 0 {
		return errors.New("tx value = 0")
	}
	if len(block.Transactions) == transactionsCountLimit && tx.Sender != storageChain {
		return errors.New("len tx = limit")
	}
	if tx.Value > startStorageReward && tx.ToStorage != storageReward {
		return errors.New("storage reward pass")
	}

	var balanceInChain uint64
	senderExpenseInTx := tx.Value + tx.ToStorage
	if value, ok := block.Mapping[tx.Sender]; ok {
		balanceInChain = value
	} else {
		balanceInChain = chain.Balance(tx.Sender)
	}
	if senderExpenseInTx > balanceInChain {
		return errors.New("expense > balance")
	}

	block.Mapping[tx.Sender] = balanceInChain - senderExpenseInTx
	block.AddBalance(chain, tx.Receiver, tx.Value)
	block.AddBalance(chain, storageChain, tx.ToStorage)
	block.Transactions = append(block.Transactions, *tx)

	return nil
}

func (block *Block) AddBalance(chain *Blockchain, receiver string, value uint64) {
	var balanceInChain uint64
	if value, ok := block.Mapping[receiver]; ok {
		balanceInChain = value
	} else {
		balanceInChain = chain.Balance(receiver)
	}
	block.Mapping[receiver] = balanceInChain + value
}

func (block *Block) Accept(chain *Blockchain, user *User, ch chan bool) error {
	if !block.transactionsIsValid(chain) {
		return errors.New("transactions is invalid")
	}
	err := block.AddTransaction(chain, &Transaction{
		RandBytes: generateRandomBytes(randBytes),
		Sender:    storageChain,
		Receiver:  user.Address(),
		Value:     storageReward,
	})
	if err != nil {
		return err
	}
	block.TimeStamp = time.Now().Format(time.RFC3339)
	block.CurrHash = block.hash()
	block.Signature = string(block.sign(user.Private()))
	block.Nonce = block.proof(ch)

	return nil
}

func (block *Block) IsValid(chain *Blockchain) bool {
	switch {
	case block == nil:
	case block.Difficulty != Difficulty:
	case !block.hashIsValid(chain, chain.Size()):
	case !block.signIsValid():
	case !block.proofIsValid():
	case !block.mappingIsValid():
	case !block.timeIsValid(chain, chain.Size()):
	case !block.transactionsIsValid(chain):
		return false
	}

	return true
}

func (block *Block) hash() []byte {
	var tempHash []byte
	for _, tx := range block.Transactions {
		tempHash = hashSum(bytes.Join(
			[][]byte{
				tempHash,
				tx.CurrHash,
			},
			[]byte{},
		))
	}

	var list []string
	for addr := range block.Mapping {
		list = append(list, addr)
	}
	sort.Strings(list)
	for _, addr := range list {
		tempHash = hashSum(bytes.Join(
			[][]byte{
				tempHash,
				[]byte(addr),
				uint64ToBytes(block.Mapping[addr]),
			},
			[]byte{},
		))
	}

	return hashSum(bytes.Join(
		[][]byte{
			tempHash,
			uint64ToBytes(uint64(Difficulty)),
			block.PrevHash,
			[]byte(block.Miner),
			[]byte(block.TimeStamp),
		},
		[]byte{},
	))
}

func (block *Block) sign(private *rsa.PrivateKey) []byte {
	return sign(private, block.CurrHash)
}

func (block *Block) proof(ch chan bool) uint64 {
	return proofOfWork(block.CurrHash, block.Difficulty, ch)
}

func (block *Block) transactionsIsValid(chain *Blockchain) bool {
	lentx := len(block.Transactions)
	plusStorage := 0
	for i := 0; i < lentx; i++ {
		if block.Transactions[i].Sender == storageChain {
			plusStorage += 1
		}

		tx := block.Transactions[i]
		if tx.Sender == storageChain {
			if tx.Receiver != block.Miner || tx.Value != storageReward {
				return false
			}
		} else {
			if !tx.hashIsValid() {
				return false
			}
			if !tx.signIsValid() {
				return false
			}
		}

		if !block.balanceIsValid(chain, tx.Sender) {
			return false
		}
		if !block.balanceIsValid(chain, tx.Receiver) {
			return false
		}

		for j := i + 1; j < lentx; j++ {
			if i < lentx-1 && bytes.Equal(block.Transactions[i].RandBytes, block.Transactions[j].RandBytes) {
				return false
			}
		}
	}

	// максимуму только одна транзакция от хранилища пользователю на блок
	if plusStorage > 1 {
		return false
	}

	if lentx == 0 || lentx > transactionsCountLimit+plusStorage {
		return false
	}

	return true
}

func (block *Block) balanceIsValid(chain *Blockchain, address string) bool {
	if _, ok := block.Mapping[address]; !ok {
		return false
	}

	lentx := len(block.Transactions)
	balanceInChain := chain.Balance(address)
	// пользователь отправил данные
	balanceSubBlock := uint64(0)
	// пользователь принял данные
	balanceAddBlock := uint64(0)
	for i := 0; i < lentx; i++ {
		tx := block.Transactions[i]
		if tx.Sender == address {
			balanceSubBlock += tx.Value + tx.ToStorage
		}
		if tx.Receiver == address {
			balanceAddBlock += tx.Value
		}
		if tx.Receiver == address && storageChain == address {
			balanceAddBlock += tx.ToStorage
		}
	}

	if (balanceInChain + balanceAddBlock - balanceSubBlock) != block.Mapping[address] {
		return false
	}

	return true
}

func (block *Block) hashIsValid(chain *Blockchain, index uint64) bool {
	if !bytes.Equal(block.hash(), block.CurrHash) {
		return false
	}

	return getIdByHash(chain.Db, Base64Encode(block.PrevHash)) == index
}

func (block *Block) signIsValid() bool {
	return verify(parsePublic(block.Miner), block.CurrHash, []byte(block.Signature)) == nil
}

func (block *Block) proofIsValid() bool {
	intHash := big.NewInt(1)
	target := big.NewInt(1)
	hash := hashSum(bytes.Join(
		[][]byte{
			block.CurrHash,
			uint64ToBytes(block.Nonce),
		},
		[]byte{},
	))
	intHash.SetBytes(hash)
	target.Lsh(target, 256-uint(block.Difficulty))

	return intHash.Cmp(target) == -1
}

func (block *Block) mappingIsValid() bool {
	for addr := range block.Mapping {
		if addr == storageChain {
			continue
		}

		flag := false
		for _, tx := range block.Transactions {
			if addr == tx.Sender || addr == tx.Receiver {
				flag = true
				break
			}
		}
		if !flag {
			return false
		}
	}

	return true
}

func (block *Block) timeIsValid(chain *Blockchain, index uint64) bool {
	btime, err := time.Parse(time.RFC3339, block.TimeStamp)
	if err != nil {
		return false
	}

	diff := time.Now().Sub(btime)
	if diff < 0 {
		return false
	}

	sblock := getBlockByHash(chain.Db, Base64Encode(block.PrevHash))
	lblock := deserializeBlock(sblock)
	if lblock == nil {
		return false
	}

	ltime, err := time.Parse(time.RFC3339, lblock.TimeStamp)
	if err != nil {
		return false
	}

	diff = btime.Sub(ltime)
	if diff > 0 {
		return false
	}

	return true
}

// New Block

const (
	Difficulty = 20 // zero byte count
)

func NewBlock(miner string, prevHash []byte) *Block {
	return &Block{
		Difficulty: Difficulty,
		PrevHash:   prevHash,
		Miner:      miner,
		Mapping:    make(map[string]uint64),
	}
}
