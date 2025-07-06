# ------------------------
# 🧪 Desenvolvimento
# ------------------------

dev-api: ## Inicia a API em modo desenvolvimento com Air
	@echo "🚀 Iniciando API em modo desenvolvimento..."
	air -c ./config/air/.air.toml

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
# 📚 Swagger Documentation
# ------------------------
.PHONY: swag-generate
swag-generate: ## 🔄 Instala e gera a documentação Swagger
	@echo "📥 Instalando swag CLI..."
	go install github.com/swaggo/swag/cmd/swag@latest
	@echo "📖 Gerando documentação Swagger..."
	swag init --dir cmd/api,internal --output docs/v1
	@echo "✅ Documentação Swagger gerada com sucesso."