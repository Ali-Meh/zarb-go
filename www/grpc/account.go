package grpc

import (
	"context"

	"github.com/zarbchain/zarb-go/account"
	"github.com/zarbchain/zarb-go/crypto"
	zarb "github.com/zarbchain/zarb-go/www/grpc/proto"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

func (zs *zarbServer) GetAccount(ctx context.Context, request *zarb.AccountRequest) (*zarb.AccountResponse, error) {
	addr, err := crypto.AddressFromString(request.Address)
	if err != nil {
		return nil, status.Errorf(codes.InvalidArgument, "Invalid address: %v", err)

	}
	acc := zs.state.Account(addr)
	if acc == nil {
		return nil, status.Errorf(codes.NotFound, "Account not found")
	}

	//attach account transactions
	trxs := make([]*zarb.Transaction, 0)
	if request.Verbosity > zarb.AccountVerbosity_ACCOUNT_INFO {
		accTrxs := zs.state.AccountTransactions(acc.Address())

		for _, trx := range accTrxs {
			trxs = append(trxs, convertTransaction(trx))
		}
	}

	res := &zarb.AccountResponse{
		Account:     convertAccount(acc),
		Tranactions: trxs,
	}

	return res, nil

}

func convertAccount(acc *account.Account) *zarb.Account {
	return &zarb.Account{
		Address:  acc.Address().String(),
		Number:   int32(acc.Number()),
		Sequence: int32(acc.Sequence()),
		Balance:  acc.Balance(),
	}
}
