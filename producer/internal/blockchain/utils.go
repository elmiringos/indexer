package blockchain

import (
	"encoding/json"
	"fmt"
	"math/big"
)

// Custom BigInt for marshalling and unmarshalling
type BigInt big.Int

func (i BigInt) MarshalJSON() ([]byte, error) {
	i2 := big.Int(i)
	return []byte(fmt.Sprintf(`"%s"`, i2.String())), nil
}

func (i *BigInt) UnmarshalJSON(data []byte) error {
	// Unmarshal as a string
	var str string
	if err := json.Unmarshal(data, &str); err != nil {
		return err
	}

	// Parse the string as a big.Int
	bi, ok := new(big.Int).SetString(str, 10)
	if !ok {
		return fmt.Errorf("failed to parse big.Int from string: %s", str)
	}

	*i = BigInt(*bi)
	return nil
}

func FromBytesToBigInt(data []byte) BigInt {
	// Create a *big.Int from the byte slice
	convertedBigInt := big.NewInt(0).SetBytes(data)

	// Convert the *big.Int to the custom BigInt type and return
	return BigInt(*convertedBigInt)
}
