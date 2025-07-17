output "postgres_port" {
  value       = var.postgres_port
  description = "The Postgres container port"
}

output "postgres_service" {
  value       = var.postgres_service
  description = "The Postgres service"
}

output "postgres_user" {
  value       = var.postgres_user
  description = "The Postgres user"
}

output "postgres_db" {
  value       = var.postgres_db
  description = "Postgres DB"
}