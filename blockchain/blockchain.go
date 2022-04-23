package blockchain

import (
	"fmt"
	"log"

	"github.com/dgraph-io/badger"
)

const (
	dbPath = "./tmp/blocks"
)

type BlockChain struct {
	LastHash []byte
	Database *badger.DB
}

type BlockChainIterator struct {
	NextHash []byte
	Database *badger.DB
}

func InitBlockChain() *BlockChain {
	var lastHash []byte

	opts := badger.DefaultOptions(dbPath)
	opts.Logger = nil

	db, err := badger.Open(opts)
	if err != nil {
		log.Panic(err)
	}

	err = db.Update(func(txn *badger.Txn) error {
		if _, err := txn.Get([]byte("lh")); err == badger.ErrKeyNotFound {
			fmt.Println("No existing blockchain found. Creating a new one...")
			genesis := Genesis()
			err = txn.Set(genesis.Hash, genesis.Serialize())
			if err != nil {
				return err
			}
			err = txn.Set([]byte("lh"), genesis.Hash)
			lastHash = genesis.Hash
			return err
		} else {
			item, err := txn.Get([]byte("lh"))
			if err != nil {
				return err
			}
			// This is the funcky way we have to read from badgerDB
			err = item.Value(func(val []byte) error {
				lastHash = append(lastHash, val...)
				return nil
			})
			return err
		}
	})

	if err != nil {
		log.Panic(err)
	}

	blockchain := BlockChain{lastHash, db}
	return &blockchain
}

func (chain *BlockChain) AddBlock(data string) {
	var lastHash []byte

	err := chain.Database.View(func(txn *badger.Txn) error {
		item, err := txn.Get([]byte("lh"))
		if err != nil {
			return err
		}
		err = item.Value(func(val []byte) error {
			lastHash = append(lastHash, val...)
			return nil
		})
		return err
	})
	if err != nil {
		log.Panic(err)
	}

	newBlock := CreateBlock(data, lastHash)

	err = chain.Database.Update(func(txn *badger.Txn) error {
		err := txn.Set(newBlock.Hash, newBlock.Serialize())
		if err != nil {
			return err
		}
		err = txn.Set([]byte("lh"), newBlock.Hash)
		chain.LastHash = newBlock.Hash
		return err
	})
	if err != nil {
		log.Panic(err)
	}
}

func (chain *BlockChain) Iterator() *BlockChainIterator {
	iter := &BlockChainIterator{chain.LastHash, chain.Database}

	return iter
}

func (iter *BlockChainIterator) Next() *Block {
	var block *Block

	err := iter.Database.View(func(txn *badger.Txn) error {
		item, err := txn.Get(iter.NextHash)
		if err != nil {
			return err
		}
		err = item.Value(func(val []byte) error {
			block = Deserialize(val)
			return nil
		})
		return err
	})
	if err != nil {
		log.Panic(err)
	}

	iter.NextHash = block.PrevHash

	return block
}
