-- block
CREATE TABLE IF NOT EXISTS "block" (
    "hash" BYTEA PRIMARY KEY,
    "number" NUMERIC NOT NULL,
    "miner_hash" BYTEA NOT NULL,
    "parent_hash" BYTEA NOT NULL,
    "gas_limit" BIGINT NOT NULL,
    "gas_used" BIGINT NOT NULL,
    "nonce" BIGINT NOT NULL,
    "size" BIGINT NOT NULL,
    "difficulty" NUMERIC NOT NULL,
    "is_pos" BOOL NOT NULL,
    "base_fee_per_gas" NUMERIC NOT NULL,
    "timestamp" BIGINT NOT NULL,
    "updated_at" TIMESTAMPTZ DEFAULT CURRENT_TIMESTAMP,
    "created_at" TIMESTAMPTZ DEFAULT CURRENT_TIMESTAMP
);

CREATE INDEX IF NOT EXISTS idx_block_number_brin ON block USING BRIN (number);
CREATE INDEX IF NOT EXISTS idx_block_timestamp_brin ON block USING BRIN (timestamp);

CREATE TRIGGER update_user_modtime
BEFORE UPDATE ON "block"
FOR EACH ROW
EXECUTE FUNCTION update_modified_column();


-- smart_contract
CREATE TABLE IF NOT EXISTS "smart_contract" (
    "address_hash" BYTEA PRIMARY KEY,
    "name" VARCHAR NOT NULL,
    "compiler_version" VARCHAR NOT NULL,
    "source_code" TEXT NOT NULL,
    "abi" JSONB NOT NULL,
    "compiler_settings" JSONB,
    "verified_by_eth" BOOL DEFAULT FALSE,
    "evm_version" VARCHAR,
    "created_at" TIMESTAMPTZ DEFAULT CURRENT_TIMESTAMP,
    "updated_at" TIMESTAMPTZ DEFAULT CURRENT_TIMESTAMP
);

CREATE INDEX IF NOT EXISTS idx_smart_contract_name ON smart_contract (name);
CREATE INDEX IF NOT EXISTS idx_smart_contract_verified ON smart_contract (verified_by_eth);

CREATE TRIGGER update_user_modtime
BEFORE UPDATE ON "smart_contract"
FOR EACH ROW
EXECUTE FUNCTION update_modified_column();


-- internal_transaction
CREATE TABLE IF NOT EXISTS "internal_transaction" (
    "block_hash" BYTEA NOT NULL,
    "index" INT NOT NULL,
    "transaction_hash" BYTEA NOT NULL,
    "status" INT NOT NULL,
    "gas" NUMERIC NOT NULL,
    "gas_used" NUMERIC NOT NULL,
    "input" BYTEA,
    "output" BYTEA,
    "amount" NUMERIC NOT NULL,
    "from_address" BYTEA NOT NULL,
    "to_address" BYTEA NOT NULL,
    "create_contract_address_hash" BYTEA NOT NULL,
    "timestamp" BIGINT NOT NULL,
    "created_at" TIMESTAMPTZ DEFAULT CURRENT_TIMESTAMP,
    "updated_at" TIMESTAMPTZ DEFAULT CURRENT_TIMESTAMP,
    PRIMARY KEY ("block_hash", "index"),
    FOREIGN KEY ("block_hash") REFERENCES "block"("hash") ON DELETE CASCADE,
    FOREIGN KEY ("create_contract_address_hash") REFERENCES "smart_contract"("address_hash") ON DELETE CASCADE
);

CREATE INDEX IF NOT EXISTS idx_internal_tx_block_hash ON internal_transaction (block_hash);
CREATE INDEX IF NOT EXISTS idx_internal_tx_timestamp_brin ON internal_transaction USING BRIN (timestamp);
CREATE INDEX IF NOT EXISTS idx_internal_tx_from_to ON internal_transaction (from_address, to_address);

CREATE TRIGGER update_user_modtime
BEFORE UPDATE ON "internal_transaction"
FOR EACH ROW
EXECUTE FUNCTION update_modified_column();


-- withdrawal
CREATE TABLE IF NOT EXISTS "withdrawal" (
    "index" INT,
    "block_hash" BYTEA,
    "address_hash" BYTEA NOT NULL,
    "validator_index" INT NOT NULL,
    "amount" NUMERIC NOT NULL,
    "created_at" TIMESTAMPTZ DEFAULT CURRENT_TIMESTAMP,
    "updated_at" TIMESTAMPTZ DEFAULT CURRENT_TIMESTAMP,
    PRIMARY KEY ("index", "block_hash"),
    FOREIGN KEY ("block_hash") REFERENCES "block"("hash") ON DELETE CASCADE
);

CREATE INDEX IF NOT EXISTS idx_withdrawal_block_hash ON withdrawal (block_hash);

CREATE TRIGGER update_user_modtime
BEFORE UPDATE ON "withdrawal"
FOR EACH ROW
EXECUTE FUNCTION update_modified_column();


-- transaction
CREATE TABLE IF NOT EXISTS "transaction" (
    "hash" BYTEA,
    "block_hash" BYTEA NOT NULL,
    "index" INT NOT NULL,
    "status" INT NOT NULL,
    "gas" NUMERIC NOT NULL,
    "gas_used" NUMERIC NOT NULL,
    "input" BYTEA,
    "value" NUMERIC NOT NULL,
    "from_address" BYTEA NOT NULL,
    "to_address" BYTEA NOT NULL,
    "nonce" BIGINT NOT NULL,
    "timestamp" BIGINT NOT NULL,
    "created_at" TIMESTAMPTZ DEFAULT CURRENT_TIMESTAMP,
    "updated_at" TIMESTAMPTZ DEFAULT CURRENT_TIMESTAMP,
    PRIMARY KEY ("hash"),
    FOREIGN KEY ("block_hash") REFERENCES "block"("hash") ON DELETE CASCADE
);

CREATE INDEX IF NOT EXISTS idx_transaction_block_hash ON transaction (block_hash);
CREATE INDEX IF NOT EXISTS idx_transaction_timestamp_brin ON transaction USING BRIN (timestamp);
CREATE INDEX IF NOT EXISTS idx_transaction_from_to ON transaction (from_address, to_address);

CREATE TRIGGER update_user_modtime
BEFORE UPDATE ON "transaction"
FOR EACH ROW
EXECUTE FUNCTION update_modified_column();


-- reward
CREATE TABLE IF NOT EXISTS "reward" (
    "block_hash" BYTEA NOT NULL,
    "address" BYTEA NOT NULL,
    "amount" NUMERIC NOT NULL,
    "created_at" TIMESTAMPTZ DEFAULT CURRENT_TIMESTAMP,
    "updated_at" TIMESTAMPTZ DEFAULT CURRENT_TIMESTAMP,
    PRIMARY KEY ("block_hash", "address"),
    FOREIGN KEY ("block_hash") REFERENCES "block"("hash") ON DELETE CASCADE
);

CREATE INDEX IF NOT EXISTS idx_reward_block_hash ON reward (block_hash);

CREATE TRIGGER update_user_modtime
BEFORE UPDATE ON "reward"
FOR EACH ROW
EXECUTE FUNCTION update_modified_column();


-- tranasction_log
CREATE TABLE IF NOT EXISTS "transaction_log" (
    "address" BYTEA NOT NULL,
    "transaction_hash" BYTEA NOT NULL,
    "block_hash" BYTEA NOT NULL,  
    "transaction_index" INT NOT NULL,
    "log_index" INT NOT NULL,
    "data" BYTEA,
    "created_at" TIMESTAMPTZ DEFAULT CURRENT_TIMESTAMP,
    "updated_at" TIMESTAMPTZ DEFAULT CURRENT_TIMESTAMP,
    PRIMARY KEY ("transaction_hash", "log_index"),
    FOREIGN KEY ("transaction_hash") REFERENCES "transaction"("hash") ON DELETE CASCADE
);

CREATE INDEX IF NOT EXISTS idx_tx_log_block_hash ON transaction_log (block_hash);
CREATE INDEX IF NOT EXISTS idx_tx_log_address ON transaction_log (address);

CREATE TRIGGER update_user_modtime
BEFORE UPDATE ON "transaction_log"
FOR EACH ROW
EXECUTE FUNCTION update_modified_column();


-- transaction_log_topic
create TABLE IF NOT EXISTS "transaction_log_topic" (
    "transaction_hash" BYTEA NOT NULL,
    "log_index" INT NOT NULL,
    "topic_index" INT NOT NULL,
    "topic" BYTEA NULL,
    "created_at" TIMESTAMPTZ DEFAULT CURRENT_TIMESTAMP,
    "updated_at" TIMESTAMPTZ DEFAULT CURRENT_TIMESTAMP,
    PRIMARY KEY ("transaction_hash", "log_index", "topic_index"),
    FOREIGN KEY ("transaction_hash", "log_index") REFERENCES "transaction_log"("transaction_hash", "log_index") ON DELETE CASCADE
);

CREATE TRIGGER update_user_modtime
BEFORE UPDATE ON "transaction_log_topic"
FOR EACH ROW
EXECUTE FUNCTION update_modified_column();


-- transaction_action
CREATE TABLE IF NOT EXISTS "transaction_action" (
    "transaction_hash" BYTEA,
    "log_index" INT,
    "data" JSONB,
    "address_contract_hash" BYTEA NOT NULL,
    "type" INT NOT NULL,
    "created_at" TIMESTAMPTZ DEFAULT CURRENT_TIMESTAMP,
    "updated_at" TIMESTAMPTZ DEFAULT CURRENT_TIMESTAMP,
    PRIMARY KEY ("transaction_hash", "log_index"),
    FOREIGN KEY ("transaction_hash", "log_index") REFERENCES "transaction_log"("transaction_hash", "log_index") ON DELETE CASCADE
);

CREATE INDEX IF NOT EXISTS idx_tx_action_address_hash ON transaction_action (address_contract_hash);

CREATE TRIGGER update_user_modtime
BEFORE UPDATE ON "transaction_action"
FOR EACH ROW
EXECUTE FUNCTION update_modified_column();


-- token
CREATE TABLE IF NOT EXISTS "token" (
    "address_hash" BYTEA PRIMARY KEY,
    "symbol" VARCHAR NOT NULL,
    "name" VARCHAR NOT NULL,
    "total_supply" NUMERIC,
    "decimals" INT NOT NULL,
    "fiat_value" NUMERIC,
    "circulation_market_cap" NUMERIC,
    "created_at" TIMESTAMPTZ DEFAULT CURRENT_TIMESTAMP,
    "updated_at" TIMESTAMPTZ DEFAULT CURRENT_TIMESTAMP
);

CREATE INDEX IF NOT EXISTS idx_token_symbol_name ON token (symbol, name);

CREATE TRIGGER update_user_modtime
BEFORE UPDATE ON "token"
FOR EACH ROW
EXECUTE FUNCTION update_modified_column();


-- token_transfer
CREATE TABLE IF NOT EXISTS "token_transfer" (
    "transaction_hash" BYTEA,
    "log_index" INT,
    "from_address" BYTEA NOT NULL,
    "to_address" BYTEA NOT NULL,
    "token_contract_address_hash" BYTEA NOT NULL,
    "amount" NUMERIC NOT NULL,
    "created_at" TIMESTAMPTZ DEFAULT CURRENT_TIMESTAMP,
    "updated_at" TIMESTAMPTZ DEFAULT CURRENT_TIMESTAMP,
    PRIMARY KEY ("transaction_hash", "log_index"),
    FOREIGN KEY ("transaction_hash") REFERENCES "transaction"("hash") ON DELETE CASCADE
);

CREATE INDEX IF NOT EXISTS idx_token_transfer_from_to ON token_transfer (from_address, to_address);
CREATE INDEX IF NOT EXISTS idx_token_transfer_token ON token_transfer (token_contract_address_hash);
CREATE INDEX IF NOT EXISTS idx_token_transfer_created_at_brin ON token_transfer USING BRIN (created_at);

CREATE TRIGGER update_user_modtime
BEFORE UPDATE ON "token_transfer"
FOR EACH ROW
EXECUTE FUNCTION update_modified_column();


-- token_instance
CREATE TABLE IF NOT EXISTS "token_instance" (
    "token_contract_address_hash" BYTEA,
    "token_id" BIGINT,
    "owner_address_hash" BYTEA NOT NULL,
    "metadata" JSONB,
    "created_at" TIMESTAMPTZ DEFAULT CURRENT_TIMESTAMP,
    "updated_at" TIMESTAMPTZ DEFAULT CURRENT_TIMESTAMP,
    PRIMARY KEY ("token_id", "token_contract_address_hash"),
    FOREIGN KEY ("token_contract_address_hash") REFERENCES "token"("address_hash") ON DELETE CASCADE
);

CREATE INDEX IF NOT EXISTS idx_token_instance_owner ON token_instance (owner_address_hash);

CREATE TRIGGER update_user_modtime
BEFORE UPDATE ON "token_instance"
FOR EACH ROW
EXECUTE FUNCTION update_modified_column();


-- audit_report
CREATE TABLE IF NOT EXISTS "audit_report" (
    "id" BIGINT PRIMARY KEY,
    "address_hash" BYTEA NOT NULL,
    "is_approved" BOOL DEFAULT FALSE,
    "submitter_name" VARCHAR NOT NULL,
    "submitter_email" VARCHAR NOT NULL,
    "audit_company_name" VARCHAR NOT NULL,
    "audit_report_url" VARCHAR NOT NULL,
    "project_url" VARCHAR NOT NULL,
    "created_at" TIMESTAMPTZ DEFAULT CURRENT_TIMESTAMP,
    "updated_at" TIMESTAMPTZ DEFAULT CURRENT_TIMESTAMP,
    FOREIGN KEY ("address_hash") REFERENCES "smart_contract"("address_hash") ON DELETE CASCADE
);

CREATE INDEX IF NOT EXISTS idx_audit_report_address_hash ON audit_report (address_hash);

CREATE TRIGGER update_user_modtime
BEFORE UPDATE ON "audit_report"
FOR EACH ROW
EXECUTE FUNCTION update_modified_column();
