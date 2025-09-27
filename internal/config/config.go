package config

import (
	"os"
	"strconv"

	"github.com/joho/godotenv"
)

const (
	DefaultPage     = "1"
	DefaultPageSize = "10"
	MaxPageSize     = 100
)

const (
	ResourceTransactions = "transactions"
	ActionCreate         = "create"
	ModelNubank          = "nubank"
)
const (
	OrderAsc  = "asc"
	OrderDesc = "desc"
)

var JwtSecret = []byte(os.Getenv("JWT_SECRET"))

type Config struct {
	DBUser string
	DBPass string
	DBHost string
	DBPort string
	DBName string

	MessageBusUser string
	MessageBusPass string
	MessageBusHost string
	MessageBusPort string

	SeedPath string
}

func LoadConfig(envPath string) (*Config, error) {
	if err := godotenv.Load(envPath); err != nil {
		return nil, err
	}

	cfg := &Config{
		DBUser: os.Getenv("DB_USER"),
		DBPass: os.Getenv("DB_PASS"),
		DBHost: os.Getenv("DB_HOST"),
		DBPort: os.Getenv("DB_PORT"),
		DBName: os.Getenv("DB_NAME"),

		MessageBusUser: os.Getenv("MESSAGE_BUS_USER"),
		MessageBusPass: os.Getenv("MESSAGE_BUS_PASS"),
		MessageBusHost: os.Getenv("MESSAGE_BUS_HOST"),
		MessageBusPort: os.Getenv("MESSAGE_BUS_PORT"),

		SeedPath: os.Getenv("SEED_PATH"),
	}

	return cfg, nil
}

type ConfigConsumer struct {
	// PollIntervalMs define o intervalo (em milissegundos) entre as tentativas de checar se a fatura foi criada no cache.
	// Utilizado enquanto o sistema aguarda outro processo criar a fatura.
	PollIntervalMs int

	// InvoiceCacheTTLMin define o tempo de vida (em minutos) que o ID da fatura ficará armazenado no cache.
	// Garante que o processo possa reutilizar esse valor em múltiplas mensagens de um mesmo "job".
	InvoiceCacheTTLMin int

	// WaitForInvoiceLimit define o tempo máximo (em segundos) que o sistema vai esperar até que o ID da fatura esteja disponível no cache.
	// Caso esse tempo seja excedido, uma falha é retornada para evitar loop infinito.
	WaitForInvoiceLimit int

	// SkipTitles define os títulos das mensagens que devem ser ignoradas pelo processamento.
	SkipTitles []string
}

func LoadConsumerConfig(envPath string) *ConfigConsumer {
	_ = godotenv.Load(envPath)

	cfg := &ConfigConsumer{
		PollIntervalMs:      getEnvAsInt("CONSUMER_POLL_INTERVAL_MS", 100),
		InvoiceCacheTTLMin:  getEnvAsInt("CONSUMER_INVOICE_CACHE_TTL_MIN", 20),
		WaitForInvoiceLimit: getEnvAsInt("CONSUMER_WAIT_FOR_INVOICE_LIMIT", 5),

		// TODO: pegar de um arquivo json
		SkipTitles: []string{"Pagamento recebido"},
	}

	return cfg
}

func getEnvAsInt(key string, defaultVal int) int {
	valStr := os.Getenv(key)
	if valStr == "" {
		return defaultVal
	}
	val, err := strconv.Atoi(valStr)
	if err != nil {
		return defaultVal
	}
	return val
}
