#!/bin/bash
set -e

# O psql -v permite injetar variáveis que o SQL entenderá como :NOME_DA_VARIAVEL
psql -v ON_ERROR_STOP=1 --username "$POSTGRES_USER" --dbname "$POSTGRES_DB" <<-EOSQL
    CREATE SCHEMA IF NOT EXISTS accounts;

    DO \$$ 
    BEGIN
        IF NOT EXISTS (SELECT FROM pg_catalog.pg_user WHERE usename = 'user_accounts') THEN
            CREATE USER user_accounts WITH PASSWORD '$ACCOUNTS_DB_PASSWORD';
        END IF;
    END \$$;

    ALTER SCHEMA accounts OWNER TO user_accounts;
    GRANT USAGE, CREATE ON SCHEMA accounts TO user_accounts;
EOSQL