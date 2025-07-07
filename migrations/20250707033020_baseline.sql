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
-- Create "debts" table
CREATE TABLE "public"."debts" (
  "id" uuid NOT NULL,
  "created_at" timestamptz NOT NULL,
  "updated_at" timestamptz NOT NULL,
  "amount" numeric(10,2) NOT NULL,
  "title" character varying NOT NULL,
  "purchase_date" timestamptz NOT NULL,
  "due_date" timestamptz NULL,
  "status" character varying NOT NULL DEFAULT 'pending',
  "category_id" uuid NULL,
  PRIMARY KEY ("id"),
  CONSTRAINT "debts_categories_category" FOREIGN KEY ("category_id") REFERENCES "public"."categories" ("id") ON UPDATE NO ACTION ON DELETE SET NULL
);
-- Create index "debt_category_id" to table: "debts"
CREATE INDEX "debt_category_id" ON "public"."debts" ("category_id");
-- Create index "debt_due_date" to table: "debts"
CREATE INDEX "debt_due_date" ON "public"."debts" ("due_date");
-- Create index "debt_due_date_category_id" to table: "debts"
CREATE INDEX "debt_due_date_category_id" ON "public"."debts" ("due_date", "category_id");
-- Create index "debt_purchase_date" to table: "debts"
CREATE INDEX "debt_purchase_date" ON "public"."debts" ("purchase_date");
-- Create index "debt_purchase_date_category_id" to table: "debts"
CREATE INDEX "debt_purchase_date_category_id" ON "public"."debts" ("purchase_date", "category_id");
