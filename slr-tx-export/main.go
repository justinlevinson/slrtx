package main

import (
	"encoding/hex"
	"fmt"
	"github.com/syndtr/goleveldb/leveldb"
	"os"
	"./slrparser"
	"strings"
)


func main() {
	f, err := os.Create("height-block-tx.csv")
		if err != nil {
			fmt.Println(err)
			return
		}
	genesisBlockHash := "e8666c8715fafbfb095132deb1dd2af63fe14d3d7163715341d48feffab458cc"
	block := dumpBlock(genesisBlockHash)
	var hash string
	var txs []slrparser.Transaction
	txs = block.Block.Transactions
		for i := 0;i<len(txs);i++ {
			fmt.Fprintln(f, block.BlockIndex.Height, block.Block.BlockHeader.Hash.String(), txs[i].Hash().String())
		}
	
	for {
		hash = block.BlockIndex.HashNext.String()
		block = dumpBlock(hash)

		
		txs = block.Block.Transactions
		
		for i := 0;i<len(txs);i++ {
			var txhash = strings.TrimSpace(txs[i].Hash().String())
			if len(txhash) > 0 {
				fmt.Fprintln(f, block.BlockIndex.Height, block.Block.BlockHeader.Hash.String(), txs[i].Hash().String())
			}
		}
		
		
	}
	
}

func dumpBlock(hashStr string) *BlockInfo{

	hash, err := hex.DecodeString(hashStr)
	if err != nil {
		panic(err)
	}

	file, err := slrparser.NewBlockFile(slrparser.SolarCoinDir() + "/blk0001.dat")
	if err != nil {
		panic(err)
	}
	defer file.Close()

	db, err := leveldb.OpenFile(slrparser.SolarCoinDir()+"/txleveldb", nil)
	if err != nil {
		panic(err)
	}
	defer db.Close()

	val, err := db.Get(append([]byte("b"), slrparser.ReverseHex(hash)...), nil)
	if err != nil {
		os.Stderr.WriteString(fmt.Sprintf("Unable to find block %v in block index.\n", hashStr))
		panic(err)
	}
	defer db.Close()

	prs := slrparser.NewBlockIndexParser(val)
	index, _ := prs.Parse()

	parser := slrparser.NewBlockParser(file, slrparser.MainnetMagicBytes)
	file.Seek(int64(index.BlockPos)-8, 0)
	block, _ := parser.ParseBlock()
	

	info := &BlockInfo{block, index}

	return(info)
}

type BlockInfo struct {
	*slrparser.Block
	*slrparser.BlockIndex
}
