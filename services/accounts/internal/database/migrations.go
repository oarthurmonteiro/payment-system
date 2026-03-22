package database

import (
	"fmt"
	"log"

	"github.com/golang-migrate/migrate/v4"
	_ "github.com/golang-migrate/migrate/v4/database/postgres"
	"github.com/golang-migrate/migrate/v4/source/iofs"
	"github.com/oarthurmonteiro/payment-system/services/accounts"
)

func RunMigrations(dbURL string) error {
    // ... (seu código de log dos arquivos) ...

    d, err := iofs.New(accounts.MigrationsFS, "migrations")
    if err != nil {
        return fmt.Errorf("erro ao ler fonte de migrations: %w", err)
    }

    // Não esqueça de fechar a fonte
    defer d.Close() 

    m, err := migrate.NewWithSourceInstance("iofs", d, dbURL)
    if err != nil {
        return fmt.Errorf("erro ao instanciar migrate: %w", err)
    }
    
    defer m.Close() 

    // log.Println("Executando m.Up()...")
    if err := m.Up(); err != nil && err != migrate.ErrNoChange {
        return fmt.Errorf("erro ao executar migrations: %w", err)
    }

    log.Println("Migrations finalizadas com sucesso!")
    return nil
}