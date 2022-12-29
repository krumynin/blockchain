package blockchain

import (
	"bytes"
	"crypto/rsa"
)

type Transaction struct {
	RandBytes []byte
	PrevBlock []byte
	Sender    string
	Receiver  string
	Value     uint64
	ToStorage uint64
	CurrHash  []byte
	Signature []byte
}

func (tx *Transaction) hash() []byte {
	return hashSum(bytes.Join(
		[][]byte{
			tx.RandBytes,
			tx.PrevBlock,
			[]byte(tx.Sender),
			[]byte(tx.Receiver),
			uint64ToBytes(tx.Value),
			uint64ToBytes(tx.ToStorage),
		},
		[]byte{},
	))
}

func (tx *Transaction) sign(private *rsa.PrivateKey) []byte {
	return sign(private, tx.CurrHash)
}

func (tx *Transaction) hashIsValid() bool {
	return bytes.Equal(tx.hash(), tx.CurrHash)
}

func (tx *Transaction) signIsValid() bool {
	return verify(parsePublic(tx.Sender), tx.CurrHash, tx.Signature) == nil
}

// NewTransaction

const (
	randBytes          = 32
	startStorageReward = 10
	storageReward      = 1
)

func NewTransaction(from *User, to string, value uint64, lastBlockHash []byte) *Transaction {
	tx := &Transaction{
		RandBytes: generateRandomBytes(randBytes),
		PrevBlock: lastBlockHash,
		Sender:    from.Address(),
		Receiver:  to,
		Value:     value,
	}
	if value > startStorageReward {
		tx.ToStorage = storageReward
	}
	tx.CurrHash = tx.hash()
	tx.Signature = tx.sign(from.Private())

	return tx
}
