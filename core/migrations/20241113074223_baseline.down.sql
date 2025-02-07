DROP TABLE IF EXISTS "audit_report";
DROP TRIGGER IF EXISTS update_session_modtime ON "audit_report";

DROP TABLE IF EXISTS "token_instance";
DROP TRIGGER IF EXISTS update_session_modtime ON "token_instance";

DROP TABLE IF EXISTS "token_transfer";
DROP TRIGGER IF EXISTS update_session_modtime ON "token_transfer";

DROP TABLE IF EXISTS "token";
DROP TRIGGER IF EXISTS update_session_modtime ON "token";

DROP TABLE IF EXISTS "transaction_action";
DROP TRIGGER IF EXISTS update_session_modtime ON "transaction_action";

DROP TABLE IF EXISTS "transaction_log";
DROP TRIGGER IF EXISTS update_session_modtime ON "transaction_log";

DROP TABLE IF EXISTS "smart_contract";
DROP TRIGGER IF EXISTS update_session_modtime ON "smart_contract";

DROP TABLE IF EXISTS "reward";
DROP TRIGGER IF EXISTS update_session_modtime ON "reward";

DROP TABLE IF EXISTS "transaction";
DROP TRIGGER IF EXISTS update_session_modtime ON "transaction";

DROP TABLE IF EXISTS "withdrawal";
DROP TRIGGER IF EXISTS update_session_modtime ON "withdrawal";

DROP TABLE IF EXISTS "internal_transaction";
DROP TRIGGER IF EXISTS update_session_modtime ON "internal_transaction";

DROP TABLE IF EXISTS "block";
DROP TRIGGER IF EXISTS update_session_modtime ON "block";