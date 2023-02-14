package blockchain

import (
	"bytes"
	"encoding/gob"
	"github.com/andrea-saboc/Building-Blockchain-in-Golang/wallet"
	"log"
)

type TxOutput struct {
	Value      int    //value in tokend
	PubKeyHash []byte //value that is needed to unlock the token inside value, very complicated scripting language called script
	//cannot dereference the part of the output
}

type TxOutputs struct {
	Outputs []TxOutput //unspent outputs
}

type TxInput struct {
	//references to the previous outputs
	ID        []byte //references the transaction the output is inside
	Out       int    //if the transaction has 3 output and if we want to reference only one of them, at index out
	Signature []byte //similiar to the pubkey in output
	PubKey    []byte
}

func NewTXOutput(value int, address string) *TxOutput {
	txo := &TxOutput{
		Value:      value,
		PubKeyHash: nil,
	}
	txo.Lock([]byte(address))

	return txo
}

func (in *TxInput) UsesKey(pubKeyHash []byte) bool {
	lockingHash := wallet.PublicKeyHash(in.PubKey)

	return bytes.Compare(lockingHash, pubKeyHash) == 0
}

func (out *TxOutput) Lock(address []byte) {
	pubKeyHash := wallet.Base58Decode(address)
	pubKeyHash = pubKeyHash[1 : len(pubKeyHash)-wallet.ChecksumLength]
	out.PubKeyHash = pubKeyHash
}

func (out *TxOutput) IsLockedWithKey(pubKeyHash []byte) bool {
	return bytes.Compare(out.PubKeyHash, pubKeyHash) == 0
}

func (outs TxOutputs) Serialize() []byte {
	var buffer bytes.Buffer
	encode := gob.NewEncoder(&buffer)
	err := encode.Encode(outs)
	if err != nil {
		log.Panic(err)
	}
	return buffer.Bytes()
}

func DeserializeOutputs(data []byte) TxOutputs {
	var outputs TxOutputs
	decode := gob.NewDecoder(bytes.NewReader(data))
	err := decode.Decode(&outputs)
	if err != nil {
		log.Panic(err)
	}
	return outputs
}
