ENV_FILE = ./config/envs/dev.env
ATLAS_CONFIG = file://./config/atlas/atlas.hcl

# Extrai as variáveis do .env e transforma em flags --var DB_USER=value
atlas_vars = $(shell grep -E '^(DB_USER|DB_PASS|DB_HOST|DB_PORT|DB_NAME|DB_DEV_NAME)=' $(ENV_FILE) | sed 's/^/--var /' | sed 's/=/=/' )


# ------------------------
# 📌 Comandos principais
# ------------------------

.PHONY: help
help: ## Lista os comandos disponíveis
	@grep -E '^[a-zA-Z_-]+:.*?## .*$$' $(MAKEFILE_LIST) | sort | awk 'BEGIN {FS = ":.*?## "}; {printf "make %-20s %s\n", $$1, $$2}'


# ------------------------
# 🧪 Desenvolvimento
# ------------------------

dev-api: ## Inicia a API em modo desenvolvimento com Air
	@echo "🚀 Iniciando API em modo desenvolvimento..."
	air -c ./config/air/.air.toml

dev-seed: ## Popula o banco com valores iniciais
	@echo "🚀 Populando banco com valores iniciais..."
	go run cmd/seed/main.go --env="./config/envs/dev.env"

dev-worker-transactions: ## Inicia o worker de transações em modo desenvolvimento com Air
	@echo "🚀 Iniciando consumer de transações em modo desenvolvimento..."
	air -c ./config/air/.air-consumer-transactions.toml

# ------------------------
# 🏗️ Ent - Codegen
# ------------------------
.PHONY: ent-generate
ent-generate: ## Gera o código Ent baseado nos schemas
	@echo "⚙️  Gerando código com Ent..."
	go get entgo.io/ent/cmd/ent@latest && \
	go run entgo.io/ent/cmd/ent generate ./internal/ent/schemas && \
	go mod tidy
	@echo "✅ Código Ent gerado com sucesso."

.PHONY: ent-clean
ent-clean: ## Remove arquivos gerados pelo Ent (exceto schemas)
	@echo "🧹 Limpando arquivos gerados pelo Ent (exceto schemas)..."
	find ./internal/ent -mindepth 1 -not -name schemas -not -path "./internal/ent/schemas/*" -exec rm -rf {} +
	@echo "✅ Limpeza concluída."

# ------------------------
# 🛠️ Atlas - Migrations
# ------------------------

.PHONY: atlas-install
atlas-install: ## Instala o Atlas CLI se não estiver instalado
	@which atlas >/dev/null || (echo "🔧 Instalando Atlas..."; curl -sSf https://atlasgo.sh | sh)

.PHONY: atlas-status
atlas-status: ## Mostra o status das migrations
	@echo "🔎 Verificando status das migrations..."
	atlas migrate status --config ${ATLAS_CONFIG} --env local $(atlas_vars)

.PHONY: atlas-up
atlas-up: ## Aplica as migrations
	@echo "⬆️  Aplicando migrations..."
	atlas migrate apply --config ${ATLAS_CONFIG} --env local $(atlas_vars)

.PHONY: atlas-down
atlas-down: ## Reverte a última migration
	@echo "↩️  Revertendo última migration..."
	atlas migrate down --config ${ATLAS_CONFIG} --env local $(atlas_vars)

.PHONY: atlas-reset
atlas-reset: ## Reverte todas as migrations
	@echo "🧨 Revertendo todas as migrations..."
	atlas migrate down --config ${ATLAS_CONFIG} --env local --all $(atlas_vars)

.PHONY: atlas-new
atlas-new: ## Cria nova migration (uso: make atlas-new NAME=descricao)
	@echo "🆕 Criando nova migration: '$(NAME)'..."
	atlas migrate diff $${NAME} --config ${ATLAS_CONFIG} --env local $(atlas_vars) --to ent://internal/ent/schemas

# ------------------------
# 📚 Swagger Documentation
# ------------------------
.PHONY: swag-generate
swag-generate: ## 🔄 Instala e gera a documentação Swagger
	@echo "📥 Instalando swag CLI..."
	go install github.com/swaggo/swag/cmd/swag@latest
	@echo "📖 Gerando documentação Swagger..."
	swag init --dir cmd/api,internal --output docs/v1
	@echo "✅ Documentação Swagger gerada com sucesso."