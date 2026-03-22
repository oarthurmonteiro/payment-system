package accounts

import "embed"

// MigrationsFS exporta os arquivos SQL para outros pacotes
//go:embed migrations/*.sql
var MigrationsFS embed.FS