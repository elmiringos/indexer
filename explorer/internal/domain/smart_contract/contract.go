package smartcontract

import (
	"context"
)

type Repository interface {
	SaveSmartContract(ctx context.Context, smartContract *SmartContract) error
}
