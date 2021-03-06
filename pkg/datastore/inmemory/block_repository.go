package inmemory

import (
	"github.com/vulcanize/vulcanizedb/pkg/core"
	"github.com/vulcanize/vulcanizedb/pkg/datastore"
)

type BlockRepository struct {
	*InMemory
}

func (blockRepository *BlockRepository) CreateOrUpdateBlock(block core.Block) (int64, error) {
	blockRepository.CreateOrUpdateBlockCallCount++
	blockRepository.blocks[block.Number] = block
	for _, transaction := range block.Transactions {
		blockRepository.receipts[transaction.Hash] = transaction.Receipt
		blockRepository.logs[transaction.TxHash] = transaction.Logs
	}
	return 0, nil
}

func (blockRepository *BlockRepository) GetBlock(blockNumber int64) (core.Block, error) {
	if block, ok := blockRepository.blocks[blockNumber]; ok {
		return block, nil
	}
	return core.Block{}, datastore.ErrBlockDoesNotExist(blockNumber)
}

func (blockRepository *BlockRepository) MissingBlockNumbers(startingBlockNumber int64, endingBlockNumber int64, nodeId string) []int64 {
	missingNumbers := []int64{}
	for blockNumber := int64(startingBlockNumber); blockNumber <= endingBlockNumber; blockNumber++ {
		if _, ok := blockRepository.blocks[blockNumber]; !ok {
			missingNumbers = append(missingNumbers, blockNumber)
		}
	}
	return missingNumbers
}

func (blockRepository *BlockRepository) SetBlocksStatus(chainHead int64) {
	for key, block := range blockRepository.blocks {
		if key < (chainHead - blocksFromHeadBeforeFinal) {
			tmp := block
			tmp.IsFinal = true
			blockRepository.blocks[key] = tmp
		}
	}
}

func (blockRepository *BlockRepository) BlockCount() int {
	return len(blockRepository.blocks)
}
