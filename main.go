package main

import (
	"crypto/sha256"
	"fmt"
	"strings"
)

// -----------------------Blockchain Data Structure:-----------------
// Block represents a block in the blockchain.
type Block struct {
	Transaction  []string
	Nonce        int
	PreviousHash string
	CurrentHash  string
	MerkleRoot   string
}

// Blockchain represents the blockchain structure.
type Blockchain struct {
	Blocks                  []Block
	NumTransactionsPerBlock int
	BlockHashMin            string
	BlockHashMax            string
}

// ------------------------Transaction Management:-----------------
// GetRecentBlock returns the most recent block in the blockchain.
func (bc *Blockchain) GetRecentBlock() Block {
	if len(bc.Blocks) == 0 {
		return Block{}
	}
	return bc.Blocks[len(bc.Blocks)-1]
}

// NewBlock creates a new block with the provided transactions.
func (bc *Blockchain) NewBlock(transactions []string) {
	if len(transactions) == 0 {
		fmt.Println("No Transactions to be added to Block.")
		return
	}

	if len(transactions) >= bc.NumTransactionsPerBlock {
		previousBlock := bc.GetRecentBlock()
		previousHash := previousBlock.CurrentHash
		merkleRoot := CreateMerkleRoot(transactions)
		nonce := FindValidNonce(merkleRoot, previousHash)
		currentHash := CreateHash(transactions, nonce, previousHash, merkleRoot)

		block := Block{
			Transaction:  transactions,
			Nonce:        nonce,
			PreviousHash: previousHash,
			CurrentHash:  currentHash,
			MerkleRoot:   merkleRoot,
		}

		bc.Blocks = append(bc.Blocks, block)
	} else {
		fmt.Println("Transactions to create a Block not reached yet.")
	}
}

// ---------------------------Merkle Tree Implementation:----------------------
// CreateMerkleRoot generates the Merkle root from a list of transactions.
func CreateMerkleRoot(transactions []string) string {
	if len(transactions) == 0 {
		return ""
	}
	if len(transactions) == 1 {
		return transactions[0]
	}

	var newTransactions []string
	for i := 0; i < len(transactions); i += 2 {
		first := transactions[i]
		second := ""
		if i+1 < len(transactions) {
			second = transactions[i+1]
		}
		combined := first + second
		newHash := sha256.Sum256([]byte(combined))
		newTransactions = append(newTransactions, fmt.Sprintf("%x", newHash))
	}

	return CreateMerkleRoot(newTransactions)
}

// CreateHash generates a hash based on transactions, nonce, previous hash, and Merkle root.
func CreateHash(transactions []string, nonce int, previousHash string, merkleRoot string) string {
	data := fmt.Sprintf("%v%d%s%s", transactions, nonce, previousHash, merkleRoot)
	hash := sha256.Sum256([]byte(data))
	return fmt.Sprintf("%x", hash)
}

// ---------------------------Proof Of work (POW) Consensus-------------------------
// FindValidNonce searches for a valid nonce to meet a specified hash prefix.
func FindValidNonce(merkleRoot string, previousHash string) int {
	nonce := 0
	leadingZeros := 4 // Adjust this number as needed

	prefix := strings.Repeat("0", leadingZeros)

	for {
		hash := CreateHash([]string{}, nonce, previousHash, merkleRoot)
		if strings.HasPrefix(hash, prefix) {
			return nonce
		}
		nonce++
	}
}

// DisplayBlocks prints the information of all blocks in the blockchain.
func (bc *Blockchain) DisplayBlocks() {
	for i, block := range bc.Blocks {
		fmt.Printf("\nBlock %d:\n", i)
		fmt.Printf("Transaction: %v\n", block.Transaction)
		fmt.Printf("Nonce: %d\n", block.Nonce)
		fmt.Printf("Previous Hash: %s\n", block.PreviousHash)
		fmt.Printf("Current Hash: %s\n", block.CurrentHash)
		fmt.Printf("Merkle Root: %s\n\n", block.MerkleRoot)
	}
}

// ChangeBlock modifies a transaction in a block at the specified index.
func (bc *Blockchain) ChangeBlock(index int, newTransaction string) bool {
	if index >= 0 && index < len(bc.Blocks) {
		oldBlock := &bc.Blocks[index]
		transactions := append([]string{}, oldBlock.Transaction...)
		transactions = append(transactions, newTransaction)

		merkleRoot := CreateMerkleRoot(transactions)
		previousHash := ""
		if index > 0 {
			previousHash = bc.Blocks[index-1].CurrentHash
		}

		nonce := FindValidNonce(merkleRoot, previousHash)
		currentHash := CreateHash(transactions, nonce, previousHash, merkleRoot)

		newBlock := &Block{
			Transaction:  transactions,
			Nonce:        nonce,
			PreviousHash: previousHash,
			CurrentHash:  currentHash,
			MerkleRoot:   merkleRoot,
		}

		bc.Blocks[index] = *newBlock
		return true
	}
	return false
}

// -----------------------------Block Validation and Consistency:-------------------------
// VerifyChain checks the validity of the blockchain by verifying the chain of blocks.
func (bc *Blockchain) VerifyChain() bool {
	for i := 1; i < len(bc.Blocks); i++ {
		if bc.Blocks[i].PreviousHash != bc.Blocks[i-1].CurrentHash {
			return false
		}
	}
	return true
}

// setNumberOfTransactionsPerBlock sets the number of transactions per block.
func (bc *Blockchain) setNumberOfTransactionsPerBlock(numTransactions int) {
	if numTransactions >= 1 {
		bc.NumTransactionsPerBlock = numTransactions
	} else {
		fmt.Println("No. of transactions per block at least be 1")
	}
}

// setBlockHashRangeForBlockCreation sets the range of block hash values for block creation.
func (bc *Blockchain) setBlockHashRangeForBlockCreation(min, max string) {
	// You can add input validation here to ensure min and max are valid hash values.
	bc.BlockHashMin = min
	bc.BlockHashMax = max
}

func main() {
	blockchain := &Blockchain{
		Blocks:                  []Block{},
		NumTransactionsPerBlock: 5, // Default value, can be changed using setNumberOfTransactionsPerBlock
	}

	blockchain.setBlockHashRangeForBlockCreation("0000", "00000") // Set your desired hash range

	var choice int
	for {
		fmt.Println("\nMenu:")
		fmt.Println("1. Set No. of Transactions per Block")
		fmt.Println("2. Set Block Hash Range")
		fmt.Println("3. Add Transactions")
		fmt.Println("4. Display Blocks")
		fmt.Println("5. Change Block Transaction")
		fmt.Println("6. Verify Blockchain")
		fmt.Println("7. Exit")
		fmt.Print("Select option: ")

		fmt.Scan(&choice)

		switch choice {
		case 1:
			var numTransactions int
			fmt.Print("Enter the number of transactions per block: ")
			fmt.Scan(&numTransactions)
			blockchain.setNumberOfTransactionsPerBlock(numTransactions)
		case 2:
			var min, max string
			fmt.Print("Enter the minimum block hash: ")
			fmt.Scan(&min)
			fmt.Print("Enter the maximum block hash: ")
			fmt.Scan(&max)
			blockchain.setBlockHashRangeForBlockCreation(min, max)
		case 3:
			var transactions []string
			for i := 0; i < blockchain.NumTransactionsPerBlock; i++ {
				fmt.Print("Enter transaction: ")
				var transaction string
				fmt.Scan(&transaction)
				transactions = append(transactions, transaction)
			}
			blockchain.NewBlock(transactions)

		case 4:
			blockchain.DisplayBlocks()
		case 5:
			var index int
			var newTransaction string
			fmt.Print("Enter index of the block to change transaction: ")
			fmt.Scan(&index)
			fmt.Print("Enter new transaction: ")
			fmt.Scan(&newTransaction)
			if blockchain.ChangeBlock(index, newTransaction) {
				fmt.Println("Block transaction changed successfully.")
			} else {
				fmt.Println("Invalid block index.")
			}
		case 6:
			if blockchain.VerifyChain() {
				fmt.Println("Blockchain is valid.")
			} else {
				fmt.Println("Blockchain is not valid.")
			}
		case 7:
			fmt.Println("Exiting...")
			return
		default:
			fmt.Println("Invalid choice. Please enter a valid option.")
		}
	}
}
