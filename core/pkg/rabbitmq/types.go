package rabbitmq

type ExchangeName string

const (
	BlockExchange               ExchangeName = "block_exchange"
	TransactionExchange         ExchangeName = "transaction_exchange"
	TokenExchange               ExchangeName = "token_exchange"
	WithdrawalExchange          ExchangeName = "withdrawal_exchange"
	TransactionLogExchange      ExchangeName = "transaction_log_exchange"
	TokenEventExchange          ExchangeName = "token_event_exchange"
	RewardExchange              ExchangeName = "reward_exchange"
	InternalTransactionExchange ExchangeName = "internal_transaction_exchange"
	TransactionActionExchange   ExchangeName = "transaction_action_exchange"
)

type RoutingKey string

const (
	BlockRoute               RoutingKey = "block_routing_key"
	TransactionRoute         RoutingKey = "transaction_routing_key"
	TokenRoute               RoutingKey = "token_routing_key"
	WithdrawalRoute          RoutingKey = "withdrawal_routing_key"
	TransactionLogRoute      RoutingKey = "transaction_log_routing_key"
	TokenEventRoute          RoutingKey = "token_event_routing_key"
	RewardRoute              RoutingKey = "reward_routing_key"
	InternalTransactionRoute RoutingKey = "internal_transaction_routing_key"
	TransactionActionRoute   RoutingKey = "transaction_action_routing_key"
)

type QueueType string

const (
	BlockQueue               QueueType = "block"
	TransactionQueue         QueueType = "transaction"
	TokenQueue               QueueType = "token"
	WithdrawalQueue          QueueType = "withdrawal"
	TransactionLogQueue      QueueType = "transaction_log"
	TokenEventQueue          QueueType = "token_event"
	RewardQueue              QueueType = "reward"
	InternalTransactionQueue QueueType = "internal_transaction"
	TransactionActionQueue   QueueType = "transaction_action"
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
