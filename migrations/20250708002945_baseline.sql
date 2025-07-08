-- Create "categories" table
CREATE TABLE "public"."categories" (
  "id" uuid NOT NULL,
  "created_at" timestamptz NOT NULL,
  "updated_at" timestamptz NOT NULL,
  "kind" character varying NOT NULL DEFAULT 'expense',
  "name" character varying NOT NULL,
  "description" character varying NULL,
  "color" character varying NULL,
  PRIMARY KEY ("id")
);
-- Create index "categories_name_key" to table: "categories"
CREATE UNIQUE INDEX "categories_name_key" ON "public"."categories" ("name");
-- Create index "category_kind" to table: "categories"
CREATE INDEX "category_kind" ON "public"."categories" ("kind");
-- Create "transactions" table
CREATE TABLE "public"."transactions" (
  "id" uuid NOT NULL,
  "created_at" timestamptz NOT NULL,
  "updated_at" timestamptz NOT NULL,
  "kind" character varying NOT NULL DEFAULT 'expense',
  "amount" numeric(10,2) NOT NULL,
  "title" character varying NOT NULL,
  "purchase_date" timestamptz NOT NULL,
  "due_date" timestamptz NULL,
  "status" character varying NOT NULL DEFAULT 'pending',
  "category_id" uuid NULL,
  PRIMARY KEY ("id"),
  CONSTRAINT "transactions_categories_category" FOREIGN KEY ("category_id") REFERENCES "public"."categories" ("id") ON UPDATE NO ACTION ON DELETE SET NULL
);
-- Create index "transaction_category_id" to table: "transactions"
CREATE INDEX "transaction_category_id" ON "public"."transactions" ("category_id");
-- Create index "transaction_due_date" to table: "transactions"
CREATE INDEX "transaction_due_date" ON "public"."transactions" ("due_date");
-- Create index "transaction_due_date_kind_category_id" to table: "transactions"
CREATE INDEX "transaction_due_date_kind_category_id" ON "public"."transactions" ("due_date", "kind", "category_id");
-- Create index "transaction_kind" to table: "transactions"
CREATE INDEX "transaction_kind" ON "public"."transactions" ("kind");
-- Create index "transaction_purchase_date" to table: "transactions"
CREATE INDEX "transaction_purchase_date" ON "public"."transactions" ("purchase_date");
-- Create index "transaction_purchase_date_kind_category_id" to table: "transactions"
CREATE INDEX "transaction_purchase_date_kind_category_id" ON "public"."transactions" ("purchase_date", "kind", "category_id");
