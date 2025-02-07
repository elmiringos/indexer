package withdrawal

type Withdrawal struct {
	Index          int
	BlockHash      string
	AddressHash    string
	ValidatorIndex int
	Amount         string
}
