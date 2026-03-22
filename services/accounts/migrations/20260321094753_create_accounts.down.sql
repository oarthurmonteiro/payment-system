-- Remover tabelas na ordem inversa das Foreign Keys
DROP TABLE IF EXISTS outbox;
DROP TABLE IF EXISTS accounts;
DROP TABLE IF EXISTS clients;

-- Remover tipos (Enums)
-- DROP TYPE IF EXISTS account_type;
DROP TYPE IF EXISTS outbox_status;
DROP TYPE IF EXISTS account_status;