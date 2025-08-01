package pg

import (
	"log"
	"os"
	"testing"

	"github.com/stretchr/testify/require"

	"CryptocoinPrice/config"
	"CryptocoinPrice/internal/app/entity"
	"CryptocoinPrice/internal/pkg/database"
)

var (
	_testRepo     *CoinRepoPG
	_testCoinUUID string

	_testCoinSymbol = "btc"
)

func TestMain(m *testing.M) {
	cfg, err := config.New()
	if err != nil {
		log.Fatalf("parse config: %v", err)
	}

	// open DB connection
	dbStorage, err := database.New(cfg.ConnString,
		database.WithTranslateError(),
		database.WithIgnoreNotFound(),
	)
	if err != nil {
		log.Fatalf("get db connection: %v", err)
	}
	_testRepo = NewCoinRepoPG(dbStorage)

	// run tests
	os.Exit(m.Run())
}

func TestCoinRepoPG_Create(t *testing.T) {
	t.Log("Create new coin")

	coin, err := _testRepo.Create(_testCoinSymbol)
	require.NoError(t, err)

	_testCoinUUID = coin.ID
	t.Logf("New coin: %+v", coin)
}

func TestCoinRepoPG_GetBySymbol(t *testing.T) {
	t.Log("Get coin by symbol")

	coin, err := _testRepo.GetBySymbol(_testCoinSymbol)
	require.NoError(t, err)

	t.Logf("Gotten coin: %+v", coin)
}

func TestCoinRepoPG_GetBySymbolUnexisting(t *testing.T) {
	t.Log("Get unexisting coin by symbol")

	_, err := _testRepo.GetBySymbol("unexisting")
	require.ErrorIs(t, err, ErrNotFound)
	t.Log("Got not found error")
}

func TestCoinRepoPG_GetObserved(t *testing.T) {
	t.Log("Get all observed coins")

	coinList, err := _testRepo.GetObserved()
	require.NoError(t, err)

	t.Logf("All observed coins: %v", coinList)
	require.Equal(t, 1, len(coinList))
}

func TestCoinRepoPG_Update(t *testing.T) {
	t.Log("Update coin")

	observed := false
	updateValues := &entity.CoinPartial{
		Observed: &observed,
	}

	err := _testRepo.Update(_testCoinUUID, updateValues)
	require.NoError(t, err)

	t.Log("Get updated coin")
	updatedCoin, err := _testRepo.GetBySymbol(_testCoinSymbol)
	require.NoError(t, err)
	t.Logf("Updated coin: %+v", updatedCoin)
}

func TestCoinRepoPG_GetObservedAfterUpdate(t *testing.T) {
	t.Log("Get all observed coins after update")

	coinList, err := _testRepo.GetObserved()
	require.NoError(t, err)

	t.Logf("All observed coins: %v", coinList)
	require.Equal(t, entity.CoinList{}, coinList)
}
