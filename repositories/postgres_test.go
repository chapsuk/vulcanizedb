package repositories_test

import (
	"github.com/8thlight/vulcanizedb/config"
	"github.com/8thlight/vulcanizedb/core"
	"github.com/8thlight/vulcanizedb/repositories"
	"github.com/jmoiron/sqlx"
	_ "github.com/lib/pq"
	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
)

var _ = Describe("The Postgres Repository", func() {

	var db *sqlx.DB
	var err error

	BeforeEach(func() {
		pgConfig := config.DbConnectionString(config.NewConfig("private").Database)
		db, err = sqlx.Connect("postgres", pgConfig)
		db.MustExec("DELETE FROM transactions")
		db.MustExec("DELETE FROM blocks")
	})

	AfterEach(func() {
		db.Close()
	})

	It("connects to the database", func() {
		Expect(err).Should(BeNil())
		Expect(db).ShouldNot(BeNil())
	})

	Describe("Saving blocks", func() {
		It("starts with no blocks", func() {
			count := repositories.NewPostgres(db).BlockCount()
			Expect(count).Should(Equal(0))
		})

		It("increments the block count", func() {
			block := core.Block{Number: 123}
			repository := repositories.NewPostgres(db)

			repository.CreateBlock(block)

			Expect(repository.BlockCount()).To(Equal(1))
		})

		It("saves the attributes of the block", func() {
			blockNumber := int64(123)
			gasLimit := int64(1000000)
			gasUsed := int64(10)
			blockHash := "x123"
			blockParentHash := "x456"
			blockNonce := "0x881db2ca900682e9a9"
			blockTime := int64(1508981640)
			uncleHash := "x789"
			blockSize := int64(1000)
			difficulty := int64(10)
			block := core.Block{
				Difficulty: difficulty,
				GasLimit:   gasLimit,
				GasUsed:    gasUsed,
				Hash:       blockHash,
				Nonce:      blockNonce,
				Number:     blockNumber,
				ParentHash: blockParentHash,
				Size:       blockSize,
				Time:       blockTime,
				UncleHash:  uncleHash,
			}

			repository := repositories.NewPostgres(db)

			repository.CreateBlock(block)

			savedBlock := repository.FindBlockByNumber(blockNumber)
			Expect(savedBlock.Difficulty).To(Equal(difficulty))
			Expect(savedBlock.GasLimit).To(Equal(gasLimit))
			Expect(savedBlock.GasUsed).To(Equal(gasUsed))
			Expect(savedBlock.Hash).To(Equal(blockHash))
			Expect(savedBlock.Nonce).To(Equal(blockNonce))
			Expect(savedBlock.Number).To(Equal(blockNumber))
			Expect(savedBlock.ParentHash).To(Equal(blockParentHash))
			Expect(savedBlock.Size).To(Equal(blockSize))
			Expect(savedBlock.Time).To(Equal(blockTime))
			Expect(savedBlock.UncleHash).To(Equal(uncleHash))
		})

		It("does not find a block when searching for a number that does not exist", func() {
			repository := repositories.NewPostgres(db)

			savedBlock := repository.FindBlockByNumber(111)

			Expect(savedBlock).To(BeNil())
		})

		It("saves one transaction associated to the block", func() {
			block := core.Block{
				Number:       123,
				Transactions: []core.Transaction{{}},
			}
			repository := repositories.NewPostgres(db)

			repository.CreateBlock(block)

			savedBlock := repository.FindBlockByNumber(123)
			Expect(len(savedBlock.Transactions)).To(Equal(1))
		})

		It("saves two transactions associated to the block", func() {
			block := core.Block{
				Number:       123,
				Transactions: []core.Transaction{{}, {}},
			}
			repository := repositories.NewPostgres(db)

			repository.CreateBlock(block)

			savedBlock := repository.FindBlockByNumber(123)
			Expect(len(savedBlock.Transactions)).To(Equal(2))
		})

		It("saves the attributes associated to a transaction", func() {
			gasLimit := int64(5000)
			gasPrice := int64(3)
			nonce := uint64(10000)
			to := "1234567890"
			value := int64(10)

			transaction := core.Transaction{
				Hash:     "x1234",
				GasPrice: gasPrice,
				GasLimit: gasLimit,
				Nonce:    nonce,
				To:       to,
				Value:    value,
			}
			block := core.Block{
				Number:       123,
				Transactions: []core.Transaction{transaction},
			}
			repository := repositories.NewPostgres(db)

			repository.CreateBlock(block)

			savedBlock := repository.FindBlockByNumber(123)
			Expect(len(savedBlock.Transactions)).To(Equal(1))
			savedTransaction := savedBlock.Transactions[0]
			Expect(savedTransaction.Hash).To(Equal(transaction.Hash))
			Expect(savedTransaction.To).To(Equal(to))
			Expect(savedTransaction.Nonce).To(Equal(nonce))
			Expect(savedTransaction.GasLimit).To(Equal(gasLimit))
			Expect(savedTransaction.GasPrice).To(Equal(gasPrice))
			Expect(savedTransaction.Value).To(Equal(value))
		})
	})

})