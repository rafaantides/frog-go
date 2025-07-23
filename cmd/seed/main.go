package main

import (
	"context"
	"encoding/json"
	"flag"
	"fmt"
	"frog-go/internal/adapters/repository/postgresql"
	"frog-go/internal/config"
	"frog-go/internal/core/domain"
	"frog-go/internal/ent/category"
	"frog-go/internal/utils"
	"frog-go/internal/utils/logger"
	"os"
)

var (
	envPath string
)

func main() {
	flag.StringVar(&envPath, "env", ".env", "Path to .env file")
	flag.Parse()
	startSeed()
}

func startSeed() {
	log := logger.NewLogger("Seed")
	ctx := context.Background()

	log.Start("🌱 Starting database seed... env: %s", envPath)

	cfg, err := config.LoadConfig(envPath)
	if err != nil {
		log.Fatal("%v", err)
	}

	repo, err := postgresql.NewPostgreSQL(
		cfg.DBUser,
		cfg.DBPass,
		cfg.DBHost,
		cfg.DBPort,
		cfg.DBName,
		cfg.SeedPath,
	)

	if err != nil {
		log.Fatal("%v", err)
	}

	postgresRepo, ok := repo.(*postgresql.PostgreSQL)
	if !ok {
		log.Fatal("❌ Type assertion failed: repo is not *postgresql.PostgreSQL")
	}

	if err := seedCategories(ctx, postgresRepo, log); err != nil {
		log.Fatal("Error seeding categories: %v", err)
	}

	if cfg.SeedPath != "" {
		if err := seedTransactions(ctx, postgresRepo, log, cfg.SeedPath); err != nil {
			log.Fatal("Error aseeding transactions: %v", err)
		}
	}

	fmt.Println("✅ Seeding completed successfully!")
}

func seedCategories(ctx context.Context, repo *postgresql.PostgreSQL, lg *logger.Logger) error {
	categories := []struct {
		Name                string
		Description         string
		Color               string
		SuggestedPercentage *int
	}{
		{
			"Assinaturas",
			"Serviços recorrentes como streaming, apps e plataformas.",
			"#FF6B6B",
			utils.IntPtr(5),
		},
		{
			"Mercado e delivery",
			"Mercado, restaurantes, delivery, cafés, padarias",
			"#FFA94D",
			utils.IntPtr(20),
		},
		{
			"Saúde e bem-estar",
			"Farmácia, plano de saúde, terapias e autocuidado.",
			"#20C997",
			utils.IntPtr(5),
		},
		{
			"Compras pessoais",
			"Produtos online, marketplaces, roupas, estética e cuidados pessoais.",
			"#845EF7",
			utils.IntPtr(10),
		},
		{
			"Transporte",
			"Uber, 99, combustível e transporte em geral.",
			"#339AF0",
			utils.IntPtr(5),
		},
		{
			"Lazer",
			"Bares, festas, eventos, shows, cinema e entretenimento.",
			"#DA77F2",
			utils.IntPtr(5),
		},
		{
			"Moradia",
			"Aluguel, condomínio, luz, água, gás e contas da casa.",
			"#FFB3C1",
			utils.IntPtr(30),
		},
		{
			"Sem categoria",
			"Gastos não classificados ou indefinidos.",
			"#CBD5E1",
			nil,
		}, {
			"Descontos",
			"Impostos, taxas e tributos diversos.",
			"#CBD5E1",
			nil,
		},
	}

	for _, c := range categories {
		exists, err := repo.Client.Category.Query().Where(category.NameEQ(c.Name)).Exist(ctx)
		if err != nil {
			return err
		}
		if exists {
			continue
		}

		_, err = repo.Client.Category.
			Create().
			SetName(c.Name).
			SetDescription(c.Description).
			SetColor(c.Color).
			SetNillableSuggestedPercentage(c.SuggestedPercentage).
			Save(ctx)
		if err != nil {
			return err
		}
		lg.Info("✅ Category created: %s", c.Name)
	}
	return nil
}

func seedTransactions(ctx context.Context, db *postgresql.PostgreSQL, lg *logger.Logger, seedPath string) error {
	data, err := os.ReadFile(seedPath + "transactions.json")
	if err != nil {
		return fmt.Errorf("erro ao ler arquivo JSON: %w", err)
	}

	var transactions []domain.Transaction
	if err := json.Unmarshal(data, &transactions); err != nil {
		return fmt.Errorf("erro ao parsear JSON: %w", err)
	}

	for _, d := range transactions {
		_, err := db.Client.Transaction.
			Create().
			SetTitle(d.Title).
			SetAmount(d.Amount).
			SetRecordDate(d.RecordDate).
			Save(ctx)
		if err != nil {
			return fmt.Errorf("erro ao criar transação '%s': %w", d.Title, err)
		}
		lg.Info("✅ Dívida criada: %s", d.Title)
	}
	return nil
}
