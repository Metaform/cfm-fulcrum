#  Copyright (c) 2025 Metaform Systems, Inc
#
#  This program and the accompanying materials are made available under the
#  terms of the Apache License, Version 2.0 which is available at
#  https://www.apache.org/licenses/LICENSE-2.0
#
#  SPDX-License-Identifier: Apache-2.0
#
#  Contributors:
#       Metaform Systems, Inc. - initial API and implementation
#

variable "namespace" {
  description = "Kubernetes namespace for deployment"
  type        = string
  default     = "default"
}

variable "replicas" {
  description = "Number of fulcrum-core replicas"
  type        = number
  default     = 1
}

variable "fulcrum_core_image" {
  description = "Docker image for fulcrum-core"
  type        = string
  default     = "fulcrum-core:latest"
}

variable "pull_policy" {
  description = "Docker image pull policy"
  type        = string
  default     = "Never"  # pull locally from Docker
}

variable "resources" {
  description = "Resource limits and requests"
  type = object({
    limits = object({
      cpu    = string
      memory = string
    })
    requests = object({
      cpu    = string
      memory = string
    })
  })
  default = {
    limits = {
      cpu    = "500m"
      memory = "512Mi"
    }
    requests = {
      cpu    = "250m"
      memory = "256Mi"
    }
  }
}

variable "labels" {
  description = "Additional labels to apply to all resources"
  type = map(string)
  default = {}
}

variable "fulcrum_core_port" {
  description = "Port that fulcrum-core HTTP server listens on"
  type        = number
  default     = 8080
}

variable "fulcrum_core_nodeport" {
  description = "NodePort HTTP server for external access"
  type        = number
  default     = 30083
}

variable "fulcrum_core_service" {
  description = "Fulcrum Core service name"
  type        = string
  default     = "fulcrum-core"
}

variable "postgres_port" {
  description = "Postgres port"
  type        = number
  default     = 5432
}

variable "postgres_user" {
  description = "Postgres user"
  type        = string
  default     = "fulcrum"
}

variable "postgres_db" {
  description = "Postgres DB"
  type        = string
  default     = "fulcrum_db"
}

variable "postgres_service" {
  description = "Postgres service"
  type        = string
  default     = "postgres-service"
}

locals {
  fulcrum_db_dsn = "host=${var.postgres_service} user=${var.postgres_user} password=password dbname=${var.postgres_db} port=${var.postgres_port} sslmode=disable"
}

variable "fulcrum_db_log_level" {
  description = "Database log level"
  type        = string
  default     = "warn"
}

variable "fulcrum_db_log_format" {
  description = "Database log format"
  type        = string
  default     = "text"
}