package blockchain

import (
	"bytes"
	"encoding/base64"
	"encoding/binary"
	"encoding/json"
)

func Base64Encode(data []byte) string {
	return base64.StdEncoding.EncodeToString(data)
}

func Base64Decode(data string) []byte {
	result, err := base64.StdEncoding.DecodeString(data)
	if err != nil {
		return nil
	}

	return result
}

func serializeBlock(block *Block) string {
	jsonData, err := json.MarshalIndent(*block, "", "\t")
	if err != nil {
		return ""
	}

	return string(jsonData)
}

func deserializeBlock(data string) *Block {
	var block Block
	err := json.Unmarshal([]byte(data), &block)
	if err != nil {
		return nil
	}

	return &block
}

func SerializeTx(tx *Transaction) string {
	jsonData, err := json.MarshalIndent(*tx, "", "\t")
	if err != nil {
		return ""
	}

	return string(jsonData)
}

func DeserializeTx(data string) *Transaction {
	var tx Transaction
	err := json.Unmarshal([]byte(data), &tx)
	if err != nil {
		return nil
	}

	return &tx
}

func uint64ToBytes(num uint64) []byte {
	var data = new(bytes.Buffer)
	err := binary.Write(data, binary.BigEndian, num)
	if err != nil {
		return nil
	}

	return data.Bytes()
}
