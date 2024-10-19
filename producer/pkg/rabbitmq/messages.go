package rabbitmq

import "math/big"

type ExchangeName string

const (
	BlockExchange           ExchangeName = "block_exchange"
	TransactionExchange     ExchangeName = "transaction_exchange"
	BlockchainEventExchange ExchangeName = "blockchain_event_exchange"
	TokenExchange           ExchangeName = "token_exchange"
	ContractExchange        ExchangeName = "contrac_exchange"
)

type QuqueType string

const (
	BlockQuque           QuqueType = "block"
	TransactionQuque     QuqueType = "transaction"
	BlockchainEventQuque QuqueType = "blockchain_event"
	TokenQuque           QuqueType = "token"
	Contractquque        QuqueType = "contract"
)

type BlockStatus int

const (
	PendingBlock BlockStatus = iota
	ConfirmedBlock
	UncleBlock
	OrphanBlock
	RevertedBlock
	FinilizedBlock
	ExecutedBlock
)

type BlockMessage struct {
	Hash       []byte      `json:"hash"`
	ParentHash []byte      `json:"parentHash"`
	Height     int64       `json:"height"`
	Status     BlockStatus `json:"status"`
	MinerHash  []byte      `json:"minerHash"`
	Difficulty uint64      `json:"difficulty"`
	Nounce     uint64      `json:"nonce"`
	GasLimit   *big.Int    `json:"gasLimit"`
	GasUsed    *big.Int    `json:"gasUsed"`
	Timestamp  int         `json:"timestamp"`
	Consensus  bool        `json:"consensus"`
	TxCount    int         `json:"TxCount"`
}

type TransactionStatus int

const (
	ConfirmTx TransactionStatus = iota
	PendingTx
	SuccessTx
	FailedTx
	ReplacedTx
)

type TransactionMessage struct {
	Hash         []byte            `json:"hash"`
	BlockHash    []byte            `json:"blockHash"`
	Index        uint64            `json:"index"`
	GasPrice     *big.Int          `json:"gasPrice"`
	GasUsed      *big.Int          `json:"GasUsed"`
	FromAddrHash []byte            `json:"fromAddrHash"`
	ToAddrHash   []byte            `json:"ToAddrHash"`
	TxStatus     TransactionStatus `json:"status"`
	TxFee        float64           `json:"txFee"`
	Value        *big.Int          `json:"value"`
	TxType       int               `json:"txType"`
	Data         []byte            `json:"data"`
	V            uint8             `json:"v"` // recovery_id
	R            *big.Int          `json:"r"` // signature_parameter
	S            *big.Int          `json:"s"` // signature_parameter
}

type BlockchainEventMessage struct {
}

type TokenMessage struct {
}

type ContractCreationMessage struct {
}

type TokenTransferMessage struct {
}

type TransactionAction struct {
}

type TokenTransfer struct {
}
