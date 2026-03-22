-- Remover tabelas na ordem inversa das Foreign Keys
DROP TABLE IF EXISTS accounts.outbox;
DROP TABLE IF EXISTS accounts.accounts;
DROP TABLE IF EXISTS accounts.clients;

-- Remover tipos (Enums)
DROP TYPE IF EXISTS accounts.account_type;
DROP TYPE IF EXISTS accounts.outbox_status;
DROP TYPE IF EXISTS accounts.account_status;