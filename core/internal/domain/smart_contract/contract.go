package smartcontract

import (
	"context"
)

type Repository interface {
	SaveSmartContract(ctx context.Context, smartContract *SmartContract) (string, error)
	SaveSmartContractAddress(ctx context.Context, addressHash string) error
	DeleteSmartContractAddress(ctx context.Context, addressHash string) error
}
