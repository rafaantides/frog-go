# 🐸 Projeto Go

Projeto Go modularizado com arquitetura hexagonal, voltado para o estudo e aplicação de boas práticas.

---

## 📂 Estrutura

```
├── cmd                     # Ponto de entrada do serviço
├── config                  # Configurações da aplicação
├── docs                    # Documentação da aplicação
├── infra                   # Arquivos de infraestrutura

├── internal                # Domínio da aplicação
│   ├── adapters            # Implementações de interfaces do core

│   ├── config              # Configurações internas

│   ├── core                # Núcleo da aplicação
│   │   ├── domain          # Entidades e tipos centrais do domínio (regras de negócio puras)
│   │   ├── dto             # Estruturas de dados para entrada e saída (Request/Response)
│   │   ├── errors          # Erros customizados da aplicação
│   │   ├── ports           # Interfaces do serviço
│   │   │   ├── inbound     # Casos de uso expostos (interfaces de serviços)
│   │   │   └── outbound    # Interfaces externas que o domínio depende (ex: repository)
│   │   └── service         # Regras de negócio (implementações dos ports/inbound)
│
│   ├── ent                 # Código gerado pelo Ent (ORM)
│       └── schemas         # Schemas Ent (mapeamento das entidades)
│   │
│   └── utils               # Funções utilitárias
│
├── migrations              # Arquivos de migration
├── static                  # Arquivos estáticos
├── go.mod
├── go.sum
├── Makefile                # Comandos úteis para desenvolvimento
└── README.md
```

---

## 🚀 Executando em desenvolvimento

### Inicia a API em modo desenvolvimento com Air

```bash
make dev-api
```

### Inicia o worker de debitos em modo desenvolvimento com Air

```bash
make dev-consumer
```

### Popula o banco com valores iniciais

```bash
make dev-seed
```

### Sobe os containers de infra em modo desenvolvimento com Docker Compose

```bash
make dev-docker-up
```

### Derruba os containers de infra em modo desenvolvimento com Docker Compose

```bash
make dev-seed
```

---

## 🧱 Migrations

### Instalar o Atlas CLI

```bash
make atlas-install
```

### Ver status das migrations

```bash
make atlas-status
```

### Aplicar migrations

```bash
make atlas-up
```

### Reverter última migration

```bash
make atlas-down
```

### Resetar banco (desfazer todas as migrations)

```bash
make atlas-reset
```

### Criar nova migration

```bash
make atlas-new NAME=descricao_da_migration
```

### Gerar snapshot atual

```bash
make atlas-snapshot
```

---

## 🧬 Gerar código Ent (ORM)

Se você alterou os schemas do Ent:

```bash
make ent-generate
```

---
