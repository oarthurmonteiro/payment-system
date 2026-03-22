-- Criar Enums dentro do schema accounts
DO $$ BEGIN
    IF NOT EXISTS (SELECT 1 FROM pg_type WHERE typname = 'account_status') THEN
        CREATE TYPE account_status AS ENUM ('PENDING', 'ACTIVE', 'BLOCKED', 'CANCELED');
    END IF;
    IF NOT EXISTS (SELECT 1 FROM pg_type WHERE typname = 'outbox_status') THEN
        CREATE TYPE outbox_status AS ENUM ('PENDING', 'PROCESSED', 'FAILED');
    END IF;
    -- IF NOT EXISTS (SELECT 1 FROM pg_type WHERE typname = 'account_type') THEN
    --     CREATE TYPE account_type AS ENUM ('checking', 'savings');
    -- END IF;
END $$;

-- Tabelas
CREATE TABLE IF NOT EXISTS clients (
    id UUID PRIMARY KEY,
    full_name TEXT NOT NULL,
    document TEXT NOT NULL,
    -- email TEXT UNIQUE NOT NULL,
    created_at TIMESTAMP WITH TIME ZONE DEFAULT NOW(),
    CONSTRAINT clients_document_unique UNIQUE (document)
);

CREATE TABLE IF NOT EXISTS accounts (
    id UUID PRIMARY KEY,
    client_id UUID REFERENCES clients(id) ON DELETE CASCADE,
    -- type account_type NOT NULL,
    status account_status DEFAULT 'ACTIVE',
    created_at TIMESTAMP WITH TIME ZONE DEFAULT NOW()
);

CREATE TABLE IF NOT EXISTS outbox (
    id UUID PRIMARY KEY,
    topic TEXT NOT NULL,
    payload JSONB NOT NULL,
    status outbox_status DEFAULT 'PENDING',
    created_at TIMESTAMP WITH TIME ZONE DEFAULT NOW()
);

-- Índices para performance
CREATE INDEX IF NOT EXISTS idx_outbox_status ON outbox(status) WHERE status = 'PENDING';
CREATE INDEX IF NOT EXISTS idx_clients_document ON clients(document);
