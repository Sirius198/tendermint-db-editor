package main

import (
	"fmt"
	"os"
	"path/filepath"

	"github.com/cosmos/cosmos-sdk/store"
	"github.com/cosmos/cosmos-sdk/store/types"
	"github.com/pkg/errors"
	dbm "github.com/tendermint/tm-db"
)

// var (
//
//	cacheSize = 100
//
// )
const (
	commitInfoKeyFmt = "s/%d" // s/<version>
)

func getCommitInfo(db dbm.DB, ver int64) (*types.CommitInfo, error) {
	cInfoKey := fmt.Sprintf(commitInfoKeyFmt, ver)

	bz, err := db.Get([]byte(cInfoKey))
	if err != nil {
		return nil, errors.Wrap(err, "failed to get commit info")
	} else if bz == nil {
		return nil, errors.New("no commit info found")
	}

	cInfo := &types.CommitInfo{}
	if err = cInfo.Unmarshal(bz); err != nil {
		return nil, errors.Wrap(err, "failed unmarshal commit info")
	}

	return cInfo, nil
}

func main() {
	fmt.Println("--- Starting ---")
	dirname, err := os.UserHomeDir()
	if err != nil {
		panic(err)
	}
	dataDir := filepath.Join(dirname, ".mun", "data")
	db, err := dbm.NewDB("application", dbm.GoLevelDBBackend, dataDir)
	if err != nil {
		panic("load db error")
	}
	cms := store.NewCommitMultiStore(db)
	ver := cms.LastCommitID().Version
	cInfo, err := getCommitInfo(db, ver)
	if err != nil {
		panic(err)
	}

	for _, storeInfo := range cInfo.StoreInfos {
		fmt.Println(storeInfo.Name, storeInfo.CommitId.Version)
	}
}
