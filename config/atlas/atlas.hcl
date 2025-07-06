variable "DB_USER" {
  type = string
}
variable "DB_PASS" {
  type = string
}
variable "DB_HOST" {
  type = string
}
variable "DB_PORT" {
  type = string
}
variable "DB_NAME" {
  type = string
}
variable "DB_DEV_NAME" {
  type = string
}

env "local" {
  url = "postgres://${var.DB_USER}:${var.DB_PASS}@${var.DB_HOST}:${var.DB_PORT}/${var.DB_NAME}?sslmode=disable"
  dev = "postgres://${var.DB_USER}:${var.DB_PASS}@${var.DB_HOST}:${var.DB_PORT}/${var.DB_DEV_NAME}?sslmode=disable"
}
