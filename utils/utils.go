package utils

import (
	"log"

	"path/filepath"

	"math/big"

	"os"

	"github.com/vulcanize/vulcanizedb/pkg/config"
	"github.com/vulcanize/vulcanizedb/pkg/core"
	"github.com/vulcanize/vulcanizedb/pkg/datastore/postgres"
	"github.com/vulcanize/vulcanizedb/pkg/geth"
)

func LoadPostgres(database config.Database, node core.Node) postgres.DB {
	db, err := postgres.NewDB(database, node)
	if err != nil {
		log.Fatalf("Error loading postgres\n%v", err)
	}
	return *db
}

func ReadAbiFile(abiFilepath string) string {
	abiFilepath = AbsFilePath(abiFilepath)
	abi, err := geth.ReadAbiFile(abiFilepath)
	if err != nil {
		log.Fatalf("Error reading ABI file at \"%s\"\n %v", abiFilepath, err)
	}
	return abi
}

func AbsFilePath(filePath string) string {
	if !filepath.IsAbs(filePath) {
		cwd, _ := os.Getwd()
		filePath = filepath.Join(cwd, filePath)
	}
	return filePath
}

func GetAbi(abiFilepath string, contractHash string, network string) string {
	var contractAbiString string
	if abiFilepath != "" {
		contractAbiString = ReadAbiFile(abiFilepath)
	} else {
		url := geth.GenURL(network)
		etherscan := geth.NewEtherScanClient(url)
		log.Printf("No ABI supplied. Retrieving ABI from Etherscan: %s", url)
		contractAbiString, _ = etherscan.GetAbi(contractHash)
	}
	_, err := geth.ParseAbi(contractAbiString)
	if err != nil {
		log.Fatalln("Invalid ABI")
	}
	return contractAbiString
}

func RequestedBlockNumber(blockNumber *int64) *big.Int {
	var _blockNumber *big.Int
	if *blockNumber == -1 {
		_blockNumber = nil
	} else {
		_blockNumber = big.NewInt(*blockNumber)
	}
	return _blockNumber
}
