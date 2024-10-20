package rabbitmq

type ExchangeName string

const (
	BlockExchange           ExchangeName = "block_exchange"
	TransactionExchange     ExchangeName = "transaction_exchange"
	BlockchainEventExchange ExchangeName = "blockchain_event_exchange"
	TokenExchange           ExchangeName = "token_exchange"
	ContractExchange        ExchangeName = "contrac_exchange"
)

type RoutingKey string

const (
	BlockRoute           RoutingKey = "block_routing_key"
	TransactionRoute     RoutingKey = "transaction_routing_key"
	BlockchainEventRoute RoutingKey = "blockchain_event_routing_key"
	TokenRoute           RoutingKey = "token_routing_key"
	Contract             RoutingKey = "contract_routing_key"
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
