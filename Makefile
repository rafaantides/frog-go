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