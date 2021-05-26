package grpc

import (
	"context"

	"github.com/zarbchain/zarb-go/crypto"
	"github.com/zarbchain/zarb-go/validator"
	zarb "github.com/zarbchain/zarb-go/www/grpc/proto"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

func (zs *zarbServer) GetValidatorByNumber(ctx context.Context, request *zarb.ValidatorByNumberRequest) (*zarb.ValidatorResponse, error) {
	validator := zs.state.ValidatorByNumber(int(request.Number))
	if validator == nil {
		return nil, status.Errorf(codes.NotFound, "NotFound Validator Address")
	}

	// TODO: make a function
	// proto validator from native validator
	return &zarb.ValidatorResponse{
		Validator: &zarb.Validator{
			PublicKey:         validator.PublicKey().String(),
			Address:           validator.Address().String(),
			Number:            int32(validator.Number()),
			Sequence:          int32(validator.Sequence()),
			Stake:             validator.Stake(),
			LastBondingHeight: int32(validator.LastBondingHeight()),
			LastJoinedHeight:  int32(validator.LastJoinedHeight()),
		},
	}, nil
}

func (zs *zarbServer) GetValidator(ctx context.Context, request *zarb.ValidatorRequest) (*zarb.ValidatorResponse, error) {
	addr, err := crypto.AddressFromString(request.Address)
	if err != nil {
		return nil, status.Errorf(codes.InvalidArgument, "Invalid Validator Address:%s", err.Error())
	}
	validator := zs.state.Validator(addr)
	if validator == nil {
		return nil, status.Errorf(codes.NotFound, "NotFound Validator Address")
	}

	return &zarb.ValidatorResponse{
		Validator: convertValidator(validator),
	}, nil
}
func (zs *zarbServer) GetValidators(ctx context.Context, request *zarb.ValidatorsRequest) (*zarb.ValidatorsResponse, error) {
	validators := zs.state.CommitteeValidators()
	validatorsResp := make([]*zarb.Validator, 0)
	for _, v := range validators {
		validatorsResp = append(validatorsResp,convertValidator(v) )
	}
	return &zarb.ValidatorsResponse{Validators: validatorsResp}, nil
}

func convertValidator(val *validator.Validator) *zarb.Validator {
	return &zarb.Validator{
		PublicKey:         val.PublicKey().String(),
		Address:           val.Address().String(),
		Number:            int32(val.Number()),
		Sequence:          int32(val.Sequence()),
		Stake:             val.Stake(),
		LastBondingHeight: int32(val.LastBondingHeight()),
		LastJoinedHeight:  int32(val.LastJoinedHeight()),
	}
}
