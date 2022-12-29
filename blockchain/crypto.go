package blockchain

import (
	"bytes"
	"crypto"
	"crypto/rand"
	"crypto/rsa"
	"crypto/sha256"
	"crypto/x509"
	"fmt"
	"math"
	"math/big"
	mrand "math/rand"
)

func generateRandomBytes(max uint) []byte {
	var slice = make([]byte, max)
	_, err := rand.Read(slice)
	if err != nil {
		return nil
	}

	return slice
}

func hashSum(data []byte) []byte {
	hash := sha256.Sum256(data)
	return hash[:]
}

func sign(private *rsa.PrivateKey, data []byte) []byte {
	signData, err := rsa.SignPSS(rand.Reader, private, crypto.SHA256, data, nil)
	if err != nil {
		return nil
	}

	return signData
}

func proofOfWork(blockHash []byte, difficulty uint8, ch chan bool) uint64 {
	var (
		target  = big.NewInt(1)
		intHash = big.NewInt(1)
		nonce   = uint64(mrand.Intn(math.MaxUint32))
		hash    []byte
	)
	target.Lsh(target, 256-uint(difficulty))

	fmt.Printf("\rNonce start: %d\n", nonce)
	for nonce < math.MaxUint64 {
		select {
		case <-ch:
			fmt.Println("lateness")
			return nonce
		default:
			hash = hashSum(bytes.Join(
				[][]byte{
					blockHash,
					uint64ToBytes(nonce),
				},
				[]byte{},
			))

			intHash.SetBytes(hash)
			if intHash.Cmp(target) == -1 {
				fmt.Printf("\rNonce end: %d\n", nonce)
				fmt.Printf("\rMining: %s\n", Base64Encode(hash))
				fmt.Println("mining success")
				return nonce
			}
		}

		nonce++
	}

	return nonce
}

func verify(public *rsa.PublicKey, data, sign []byte) error {
	return rsa.VerifyPSS(public, crypto.SHA256, data, sign, nil)
}

func generatePrivate(bits uint) *rsa.PrivateKey {
	private, err := rsa.GenerateKey(rand.Reader, int(bits))
	if err != nil {
		return nil
	}

	return private
}

func stringPrivate(private *rsa.PrivateKey) string {
	return Base64Encode(x509.MarshalPKCS1PrivateKey(private))
}

func parsePrivate(privateData string) *rsa.PrivateKey {
	pub, err := x509.ParsePKCS1PrivateKey(Base64Decode(privateData))
	if err != nil {
		return nil
	}

	return pub
}

func stringPublic(public *rsa.PublicKey) string {
	return Base64Encode(x509.MarshalPKCS1PublicKey(public))
}

func parsePublic(publicData string) *rsa.PublicKey {
	pub, err := x509.ParsePKCS1PublicKey(Base64Decode(publicData))
	if err != nil {
		return nil
	}

	return pub
}
