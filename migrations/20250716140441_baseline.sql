-- Create "categories" table
CREATE TABLE "public"."categories" (
  "id" uuid NOT NULL,
  "created_at" timestamptz NOT NULL,
  "updated_at" timestamptz NOT NULL,
  "name" character varying NOT NULL,
  "description" character varying NULL,
  "color" character varying NULL,
  PRIMARY KEY ("id")
);
-- Create index "categories_name_key" to table: "categories"
CREATE UNIQUE INDEX "categories_name_key" ON "public"."categories" ("name");
-- Create "invoices" table
CREATE TABLE "public"."invoices" (
  "id" uuid NOT NULL,
  "created_at" timestamptz NOT NULL,
  "updated_at" timestamptz NOT NULL,
  "status" character varying NOT NULL DEFAULT 'pending',
  "amount" numeric(10,2) NOT NULL DEFAULT 0,
  "title" character varying NOT NULL,
  "due_date" timestamptz NOT NULL,
  PRIMARY KEY ("id")
);
-- Create "transactions" table
CREATE TABLE "public"."transactions" (
  "id" uuid NOT NULL,
  "created_at" timestamptz NOT NULL,
  "updated_at" timestamptz NOT NULL,
  "record_type" character varying NOT NULL DEFAULT 'expense',
  "status" character varying NOT NULL DEFAULT 'pending',
  "amount" numeric(10,2) NOT NULL,
  "title" character varying NOT NULL,
  "record_date" timestamptz NOT NULL,
  "invoice_id" uuid NULL,
  "category_id" uuid NULL,
  PRIMARY KEY ("id"),
  CONSTRAINT "transactions_categories_category" FOREIGN KEY ("category_id") REFERENCES "public"."categories" ("id") ON UPDATE NO ACTION ON DELETE SET NULL,
  CONSTRAINT "transactions_invoices_invoice" FOREIGN KEY ("invoice_id") REFERENCES "public"."invoices" ("id") ON UPDATE NO ACTION ON DELETE CASCADE
);
-- Create index "transaction_category_id" to table: "transactions"
CREATE INDEX "transaction_category_id" ON "public"."transactions" ("category_id");
-- Create index "transaction_invoice_id" to table: "transactions"
CREATE INDEX "transaction_invoice_id" ON "public"."transactions" ("invoice_id");
-- Create index "transaction_record_date" to table: "transactions"
CREATE INDEX "transaction_record_date" ON "public"."transactions" ("record_date");
-- Create index "transaction_record_date_record_type_category_id" to table: "transactions"
CREATE INDEX "transaction_record_date_record_type_category_id" ON "public"."transactions" ("record_date", "record_type", "category_id");
-- Create index "transaction_record_type" to table: "transactions"
CREATE INDEX "transaction_record_type" ON "public"."transactions" ("record_type");
