# Payment System

A project focused on reinforcing concepts and practicing system design strategies.
Goal is develop microsservices and a ecosystem that works even if something fail.

Stack:
- **Postgres:** The phisical storage database, but will have different schemas and different users per service;
- **Public API Service(Go):** Receives the communication with the client and orchestrate the requests;
- **Accounts Service (?):** Manage client and accounts information;
- **Payments Service (C#):** Manage payments approvals and rejections;
- **Ledger Service (Go):** Manage the ledger and accounts balance information;
- **PSP (Bun):** Simulate an external PSP (Payment Service Provider) that can succed, fail, delay responses;
- **Kafka:** Event Storage
- **Docker:** Containerize everything and split in different networks to simulate communication between servers 


```sh
# Gera a senha do Admin do Banco
echo "POSTGRES_PASSWORD=$(openssl rand -hex 32)" >> .env
# Gera a senha específica do Accounts Service
echo "ACCOUNTS_DB_PASSWORD=$(openssl rand -hex 32)" >> .env
```

```
payments-system/
└── contracts/
    └── accounts/     <-- Esta é a sua "root" para o Buf
        └── v1/
            └── onboarding.proto
```