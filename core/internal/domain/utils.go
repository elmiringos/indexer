package domain

import (
	"database/sql/driver"
	"encoding/hex"
	"encoding/json"
	"errors"
	"fmt"
	"math/big"
	"strings"
)

var (
	ErrInvalidHexString = errors.New("invalid hex string")
)

// Custom BigInt for marshalling and unmarshalling
type BigInt big.Int

func (i *BigInt) Bytes() []byte {
	return (*big.Int)(i).Bytes()
}

func (i BigInt) String() string {
	return (*big.Int)(&i).String()
}

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

func (i BigInt) Value() (driver.Value, error) {
	return i.String(), nil
}

// Convert SQL string to big.Int
func (i *BigInt) Scan(value interface{}) error {
	str, ok := value.(string)
	if !ok {
		return errors.New("failed to scan BigInt")
	}
	bi, ok := new(big.Int).SetString(str, 10)
	if !ok {
		return errors.New("invalid BigInt format")
	}
	*i = BigInt(*bi)
	return nil
}

func FromBytesToBigInt(data []byte) *BigInt {
	return (*BigInt)(big.NewInt(0).SetBytes(data))
}

func HexFromBinary(data []byte) string {
	return "0x" + hex.EncodeToString(data)
}

func HexToBinary(hexString string) ([]byte, error) {
	if !strings.HasPrefix(hexString, "0x") {
		return nil, ErrInvalidHexString
	}
	return hex.DecodeString(hexString[2:])
}
