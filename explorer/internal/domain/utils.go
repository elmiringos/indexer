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

type HexBytes []byte

func (b *HexBytes) UnmarshalJSON(input []byte) error {
	// Remove quotes
	strInput := strings.Trim(string(input), "\"")

	// Remove 0x if present
	strInput = strings.TrimPrefix(strInput, "0x")

	data, err := hex.DecodeString(strInput)
	if err != nil {
		return err
	}

	*b = data
	return nil
}

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

// Helper to convert *BigInt to *big.Int
func (i *BigInt) toBigInt() *big.Int {
	return (*big.Int)(i)
}

// Helper to wrap big.Int back into BigInt
func fromBigInt(i *big.Int) BigInt {
	return *(*BigInt)(i)
}

// Cmp compares i with j
func (i BigInt) Cmp(j BigInt) int {
	return (*big.Int)(&i).Cmp((*big.Int)(&j))
}

// Sum adds two BigInt values and returns a new BigInt
func (i BigInt) Sum(j BigInt) BigInt {
	result := new(big.Int).Add(i.toBigInt(), j.toBigInt())
	return fromBigInt(result)
}

// Div divides i by j and returns a new BigInt
func (i BigInt) Sub(j BigInt) BigInt {
	result := new(big.Int).Sub(i.toBigInt(), j.toBigInt())
	return fromBigInt(result)
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

func BigIntZero() BigInt {
	return BigInt(*big.NewInt(0))
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
