package every_block

import (
	"fmt"
	"github.com/vulcanize/vulcanizedb/pkg/datastore/postgres"
	"log"
)

type ERC20RepositoryInterface interface {
	Create(supply TokenSupply) error
	MissingBlocks(startingBlock int64, highestBlock int64) ([]int64, error)
}

type TokenSupplyRepository struct {
	*postgres.DB
}

type repositoryError struct {
	err         string
	msg         string
	blockNumber int64
}

func (re *repositoryError) Error() string {
	return fmt.Sprintf(re.msg, re.blockNumber, re.err)
}

func newRepositoryError(err error, msg string, blockNumber int64) error {
	e := repositoryError{err.Error(), msg, blockNumber}
	log.Println(e.Error())
	return &e
}

const (
	GetBlockError          = "Error fetching block number %d: %s"
	InsertTokenSupplyError = "Error inserting token_supply for block number %d: %s"
	MissingBlockError      = "Error finding missing token_supply records starting at block %d: %s"
)

func (tsp *TokenSupplyRepository) Create(supply TokenSupply) error {
	var blockId int
	err := tsp.DB.Get(&blockId, `SELECT id FROM blocks WHERE number = $1`, supply.BlockNumber)
	if err != nil {
		return newRepositoryError(err, GetBlockError, supply.BlockNumber)
	}

	_, err = tsp.DB.Exec(
		`INSERT INTO token_supply (supply, token_address, block_id)
                VALUES($1, $2, $3)`,
		supply.Value, supply.TokenAddress, blockId)
	if err != nil {
		return newRepositoryError(err, InsertTokenSupplyError, supply.BlockNumber)
	}
	return nil
}

func (tsp *TokenSupplyRepository) MissingBlocks(startingBlock int64, highestBlock int64) ([]int64, error) {
	blockNumbers := make([]int64, 0)

	err := tsp.DB.Select(
		&blockNumbers,
		`SELECT number FROM BLOCKS
               LEFT JOIN token_supply ON blocks.id = block_id
               WHERE block_id ISNULL
               AND number >= $1
               AND number <= $2
               LIMIT 20`,
		startingBlock,
		highestBlock,
	)
	if err != nil {
		return []int64{}, newRepositoryError(err, MissingBlockError, startingBlock)
	}
	return blockNumbers, err
}