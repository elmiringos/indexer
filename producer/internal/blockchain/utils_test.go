package blockchain_test

import (
	"encoding/json"
	"math/big"
	"testing"

	"github.com/elmiringos/indexer/producer/internal/blockchain"
	"github.com/stretchr/testify/assert"
)

func TestBigInt_String(t *testing.T) {
	// Test that BigInt correctly returns the string representation of the big.Int
	bi := blockchain.BigInt(*big.NewInt(123456789))
	assert.Equal(t, "123456789", bi.String())
}

func TestBigInt_MarshalJSON(t *testing.T) {
	// Test that BigInt is marshaled to a JSON string
	bi := blockchain.BigInt(*big.NewInt(123456789))
	expectedJSON := `"123456789"`

	// Marshal the BigInt into JSON
	result, err := json.Marshal(bi)
	assert.NoError(t, err)
	assert.Equal(t, expectedJSON, string(result))
}

func TestBigInt_UnmarshalJSON(t *testing.T) {
	// Test that JSON string is unmarshaled correctly into BigInt
	jsonData := `"123456789"`
	var bi blockchain.BigInt

	// Unmarshal the JSON into BigInt
	err := json.Unmarshal([]byte(jsonData), &bi)
	assert.NoError(t, err)

	// Verify the result
	expectedBigInt := blockchain.BigInt(*big.NewInt(123456789))
	assert.Equal(t, expectedBigInt, bi)
}

func TestBigInt_UnmarshalJSON_Error(t *testing.T) {
	// Test that invalid JSON (non-numeric) returns an error
	invalidJSONData := `"not-a-number"`
	var bi blockchain.BigInt

	// Try unmarshaling invalid JSON into BigInt
	err := json.Unmarshal([]byte(invalidJSONData), &bi)
	assert.Error(t, err)
}

func TestFromBytesToBigInt(t *testing.T) {
	// Test that FromBytesToBigInt converts a byte slice to BigInt
	data := []byte{0x01, 0x02, 0x03, 0x04}
	expectedBigInt := blockchain.BigInt(*big.NewInt(16909060)) // big.NewInt(16909060) is the integer representation of the byte slice

	// Convert the byte slice to BigInt
	result := blockchain.FromBytesToBigInt(data)
	assert.Equal(t, expectedBigInt, result)
}
