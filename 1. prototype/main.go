package main

/*
블록체인 기본 프로토타입
*/

import (
	"bytes"
	"crypto/sha256"
	"fmt"
	"math"
	"math/big"
	"strconv"
	"time"
)

// 난이도
const targetBits = 24

// 배열을 활용하여 블록체인 구조 구현
type Blcokchain struct {
	blocks []*Block
}

type ProofOfWork struct {
	block  *Block
	target *big.Int
}

// 블록 구성
type Block struct {
	Timestamp     int64  // 블록 생성 시간
	Data          []byte // 블록에 포한된 정보
	PrevBlockHash []byte // 이전 블록의 해시값
	Hash          []byte // 해당 블록의 해시값
	Nonce         int
}

func NewProofOfWork(b *Block) *ProofOfWork {
	target := big.NewInt(1)
	target.Lsh(target, uint(256-targetBits)) // 256 - targetBits Left Shift
	// target : 0x10000000000000000000000000000000000000000000000000000000000

	pow := &ProofOfWork{b, target}
	return pow
}

func (pow *ProofOfWork) Run() (int, []byte) {
	var hashInt big.Int
	var hash [32]byte
	nonce := 0

	fmt.Printf("Mining the block containing \"%s\"\n", pow.block.Data)

	for nonce < math.MaxInt64 {
		data := pow.prepareData(nonce)
		hash = sha256.Sum256(data)
		fmt.Printf("\r%x", hash)

		hashInt.SetBytes(hash[:])
		if hashInt.Cmp(pow.target) == -1 {
			break
		} else {
			nonce++
		}
	}
	fmt.Print("\n\n")

	return nonce, hash[:]
}

func IntToHex(n int64) []byte {
	return []byte(strconv.FormatInt(n, 16))
}

func (pow *ProofOfWork) prepareData(nonce int) []byte {
	data := bytes.Join(
		[][]byte{
			pow.block.PrevBlockHash,
			pow.block.Data,
			IntToHex(pow.block.Timestamp),
			IntToHex(int64(targetBits)),
			IntToHex(int64(nonce)),
		},
		[]byte{},
	)
	return data
}

func (pow *ProofOfWork) Validate() bool {
	var hashInt big.Int

	data := pow.prepareData(pow.block.Nonce)
	hash := sha256.Sum256(data)
	hashInt.SetBytes(hash[:])

	isValid := hashInt.Cmp(pow.target) == -1
	return isValid
}

// 블록을 구성하는 필드들을 하나로 이은 뒤 SHA-256으로 해시화
func (b *Block) SetHash() {
	timestamp := []byte(strconv.FormatInt(b.Timestamp, 10))
	headers := bytes.Join([][]byte{b.PrevBlockHash, b.Data, timestamp}, []byte{})
	hash := sha256.Sum256(headers)
	b.Hash = hash[:]
}

// 블록 생성
func NewBlock(data string, prevBlockHash []byte) *Block {
	block := &Block{time.Now().Unix(), []byte(data), prevBlockHash, []byte{}, 0}
	pow := NewProofOfWork(block)
	nonce, hash := pow.Run()

	// block.SetHash()
	block.Hash = hash[:]
	block.Nonce = nonce
	return block
}

// 블록 추가
func (bc *Blcokchain) AddBlock(data string) {
	prevBlock := bc.blocks[len(bc.blocks)-1]
	newBlock := NewBlock(data, prevBlock.Hash)
	bc.blocks = append(bc.blocks, newBlock)
}

// 처음엔 적어도 하나의 블록이 필요 → 초기 블록 : 제네시스 블록
func NewGenesisBlock() *Block {
	return NewBlock("Genesis Block", []byte{})
}

// 제네시스 블록을 가지고 블록체인 생성
func NewBlockchain() *Blcokchain {
	return &Blcokchain{[]*Block{NewGenesisBlock()}}
}

func main() {
	bc := NewBlockchain()
	bc.AddBlock("Send 1 BTC to Ivan")
	bc.AddBlock("Send 2 more BTC to Ivan")

	for _, block := range bc.blocks {
		fmt.Printf("Prev. hash: %x\n", block.PrevBlockHash)
		fmt.Printf("Data: %s\n", block.Data)
		fmt.Printf("Hash: %x\n\n", block.Hash)

		pow := NewProofOfWork(block)
		fmt.Printf("PoW: %s\n", strconv.FormatBool(pow.Validate()))
		fmt.Println()
	}
}

/*
실제 블록체인에서 새로운 데이터를 추가하는데는 몇 가지
작업(권한을 얻기 위한 작업 증명(Proof of Work) 등)이 더 필요하다 !

하나의 새로운 블록은 반드시 네트워크의 참여자들로부터 확인과 승인을 받아야한다 !
*/
