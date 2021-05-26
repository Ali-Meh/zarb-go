package grpc

import (
	"fmt"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/zarbchain/zarb-go/account"
	"github.com/zarbchain/zarb-go/tx"
	"github.com/zarbchain/zarb-go/util"
	zarb "github.com/zarbchain/zarb-go/www/grpc/proto"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

func TestGetAccount(t *testing.T) {
	conn, client := callServer(t)
	defer conn.Close()

	acc1, sig1 := account.GenerateTestAccount(util.RandInt(10000))
	t.Run("Should return nil for non existing account ", func(t *testing.T) {
		res, err := client.GetAccount(tCtx, &zarb.AccountRequest{Address: acc1.Address().String()})
		assert.Error(t, err)
		assert.Equal(t, err.Error(), status.Errorf(codes.NotFound, "Account not found").Error())
		assert.Nil(t, res)
	})

	tMockState.Store.UpdateAccount(acc1)
	t.Run("Should only return account details with verbosity 0", func(t *testing.T) {
		res, err := client.GetAccount(tCtx, &zarb.AccountRequest{Address: acc1.Address().String(), Verbosity: zarb.AccountVerbosity_ACCOUNT_TRANSACTIONS})
		fmt.Println("the vlaue in", res)
		assert.NoError(t, err)
		assert.NotNil(t, res)
		assert.Equal(t, acc1.Balance(), res.Account.Balance)
		assert.Empty(t, res.Tranactions)
	})

	acc2, _ := account.GenerateTestAccount(util.RandInt(10000))
	trx1 := tx.NewSendTx(tMockState.GenesisHash(), 1, acc1.Address(), acc2.Address(), 1, 1000, "")
	sig1.SignMsg(trx1)

	tMockState.Store.SaveTransaction(trx1)

	t.Run("Should only return account details with verbosity 1", func(t *testing.T) {
		res, err := client.GetAccount(tCtx, &zarb.AccountRequest{Address: acc1.Address().String(), Verbosity: zarb.AccountVerbosity_ACCOUNT_TRANSACTIONS})
		fmt.Println("the vlaue in", res)
		assert.NoError(t, err)
		assert.NotNil(t, res)
		assert.Equal(t, acc1.Balance(), res.Account.Balance)
		assert.NotEmpty(t, res.Tranactions)
		assert.Equal(t, trx1.ID().String(), res.Tranactions[0].Id)
	})
}
