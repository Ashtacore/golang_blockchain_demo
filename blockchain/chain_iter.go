package blockchain

import (
	"log"

	"github.com/dgraph-io/badger"
)

type BlockChainIterator struct {
	NextHash []byte
	Database *badger.DB
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
