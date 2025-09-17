-- Enable UUID extension
CREATE EXTENSION IF NOT EXISTS "uuid-ossp";

-- Custom types
CREATE TYPE reward_type AS ENUM ('ONBOARDING', 'REFERRAL', 'TRADING_MILESTONE', 'BONUS');
CREATE TYPE account_type AS ENUM ('STOCK_ASSET', 'CASH_OUTFLOW', 'FEES_EXPENSE');
CREATE TYPE fee_type AS ENUM ('BROKERAGE', 'STT', 'GST', 'STAMP_DUTY', 'SEBI_FEE');
CREATE TYPE action_type AS ENUM ('SPLIT', 'MERGER', 'DELISTING', 'BONUS_ISSUE');
CREATE TYPE adjustment_type AS ENUM ('CORRECTION', 'REFUND', 'BONUS', 'PENALTY');

-- Users
CREATE TABLE IF NOT EXISTS users (
    user_id VARCHAR(50) PRIMARY KEY,
    email VARCHAR(255) UNIQUE,
    created_at TIMESTAMPTZ DEFAULT NOW(),
    updated_at TIMESTAMPTZ DEFAULT NOW()
);

-- Stock symbols
CREATE TABLE IF NOT EXISTS stock_symbols (
    symbol VARCHAR(20) PRIMARY KEY,
    company_name VARCHAR(200) NOT NULL,
    exchange VARCHAR(10) NOT NULL CHECK (exchange IN ('NSE', 'BSE')),
    isin VARCHAR(12),
    is_active BOOLEAN DEFAULT TRUE,
    created_at TIMESTAMPTZ DEFAULT NOW(),
    updated_at TIMESTAMPTZ DEFAULT NOW()
);


-- Stock rewards (partitioned by reward_timestamp month)
CREATE TABLE IF NOT EXISTS stock_rewards (
    reward_id UUID PRIMARY KEY DEFAULT uuid_generate_v4(),
    user_id VARCHAR(50) NOT NULL REFERENCES users(user_id),
    stock_symbol VARCHAR(20) NOT NULL REFERENCES stock_symbols(symbol),
    shares DECIMAL(18,6) NOT NULL CHECK (shares > 0),
    reward_type reward_type NOT NULL,
    price_at_reward DECIMAL(18,4) NOT NULL CHECK (price_at_reward > 0),
    total_stock_value DECIMAL(18,4) NOT NULL CHECK (total_stock_value > 0),
    idempotency_key VARCHAR(100) UNIQUE NOT NULL,
    reward_timestamp TIMESTAMPTZ NOT NULL,
    created_at TIMESTAMPTZ DEFAULT NOW(),
    CONSTRAINT chk_value_calculation CHECK (
        ABS(total_stock_value - (shares * price_at_reward)) < 0.01
    )
) PARTITION BY RANGE (reward_timestamp);

CREATE INDEX IF NOT EXISTS idx_stock_rewards_user_date ON stock_rewards (user_id, DATE(reward_timestamp));
CREATE INDEX IF NOT EXISTS idx_stock_rewards_timestamp ON stock_rewards (reward_timestamp);
CREATE INDEX IF NOT EXISTS idx_stock_rewards_symbol ON stock_rewards (stock_symbol);

-- Example partition for 2025-09
DO $$
BEGIN
    IF NOT EXISTS (SELECT 1 FROM pg_class WHERE relname = 'stock_rewards_y2025m09') THEN
        CREATE TABLE stock_rewards_y2025m09 PARTITION OF stock_rewards
            FOR VALUES FROM ('2025-09-01') TO ('2025-10-01');
    END IF;
END $$;

-- Ledger entries
CREATE TABLE IF NOT EXISTS ledger_entries (
    entry_id UUID PRIMARY KEY DEFAULT uuid_generate_v4(),
    transaction_id UUID NOT NULL,
    user_id VARCHAR(50) NOT NULL REFERENCES users(user_id),
    account_type account_type NOT NULL,
    stock_symbol VARCHAR(20) REFERENCES stock_symbols(symbol),
    debit_amount DECIMAL(18,4) DEFAULT 0 CHECK (debit_amount >= 0),
    credit_amount DECIMAL(18,4) DEFAULT 0 CHECK (credit_amount >= 0),
    shares DECIMAL(18,6) DEFAULT 0,
    description TEXT,
    reference_id UUID,
    entry_timestamp TIMESTAMPTZ NOT NULL,
    created_at TIMESTAMPTZ DEFAULT NOW(),
    CONSTRAINT chk_debit_or_credit CHECK (
        (debit_amount > 0 AND credit_amount = 0) OR 
        (credit_amount > 0 AND debit_amount = 0)
    )
);

CREATE INDEX IF NOT EXISTS idx_ledger_user_account ON ledger_entries (user_id, account_type);
CREATE INDEX IF NOT EXISTS idx_ledger_transaction ON ledger_entries (transaction_id);
CREATE INDEX IF NOT EXISTS idx_ledger_reference ON ledger_entries (reference_id);

-- Reward fees
CREATE TABLE IF NOT EXISTS reward_fees (
    fee_id UUID PRIMARY KEY DEFAULT uuid_generate_v4(),
    reward_id UUID NOT NULL REFERENCES stock_rewards(reward_id) ON DELETE CASCADE,
    fee_type fee_type NOT NULL,
    amount DECIMAL(18,4) NOT NULL CHECK (amount >= 0),
    rate DECIMAL(8,6),
    calculation_basis VARCHAR(50),
    created_at TIMESTAMPTZ DEFAULT NOW()
);

-- Stock prices
CREATE TABLE IF NOT EXISTS stock_prices (
    price_id UUID PRIMARY KEY DEFAULT uuid_generate_v4(),
    stock_symbol VARCHAR(20) NOT NULL REFERENCES stock_symbols(symbol),
    price DECIMAL(18,4) NOT NULL CHECK (price > 0),
    price_timestamp TIMESTAMPTZ NOT NULL,
    source VARCHAR(50) DEFAULT 'API',
    created_at TIMESTAMPTZ DEFAULT NOW()
);

CREATE UNIQUE INDEX IF NOT EXISTS idx_stock_prices_symbol_timestamp ON stock_prices (stock_symbol, price_timestamp);
CREATE INDEX IF NOT EXISTS idx_stock_prices_symbol_time_desc ON stock_prices (stock_symbol, price_timestamp DESC);

-- Portfolio snapshots
CREATE TABLE IF NOT EXISTS portfolio_snapshots (
    snapshot_id UUID PRIMARY KEY DEFAULT uuid_generate_v4(),
    user_id VARCHAR(50) NOT NULL REFERENCES users(user_id),
    stock_symbol VARCHAR(20) NOT NULL REFERENCES stock_symbols(symbol),
    total_shares DECIMAL(18,6) NOT NULL,
    price DECIMAL(18,4) NOT NULL,
    total_value DECIMAL(18,4) NOT NULL,
    snapshot_date DATE NOT NULL,
    created_at TIMESTAMPTZ DEFAULT NOW()
);

CREATE UNIQUE INDEX IF NOT EXISTS idx_portfolio_snapshots_user_stock_date ON portfolio_snapshots (user_id, stock_symbol, snapshot_date);


