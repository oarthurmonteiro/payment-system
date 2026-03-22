```
internal/
├── domain/
│   ├── account.go      # Struct + Interface AccountRepository
│   ├── client.go       # Struct + Interface ClientRepository
|   ├── document.go     # Value Object CPF
|   └── onboarding.go   # OnBoarding Interface
├── usecase/
│   └── onboarding.go   # Orquestração (Application Service)
└── database/           # IMPLEMENTAÇÃO (Infra)
    ├── postgres.go
    └── client_repo.go  # Aqui ele implementa a interface do domain
```

```sh
# 1. Instalar o compilador protoc (se não tiver)
# macOS: brew install protobuf
# Ubuntu: sudo apt install -y protobuf-compiler

# 2. Instalar os plugins de Go
go install google.golang.org/protobuf/cmd/protoc-gen-go@latest
go install google.golang.org/grpc/cmd/protoc-gen-go-grpc@latest
```

`proto-path` está usando um caminho relativo pois estamos em um monorepo
```sh
# Criamos a pasta gen se não existir
mkdir -p gen

protoc --proto_path=../../contracts \
    --go_out=./gen --go_opt=paths=source_relative \
    --go-grpc_out=./gen --go-grpc_opt=paths=source_relative \
    ../../contracts/accounts/v1/onboarding.proto
```

```sh
go test -coverprofile=cover.out ./...
go tool cover -html=cover.out
```