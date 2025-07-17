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
  description = "Number of agent replicas"
  type        = number
  default     = 1
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
      cpu    = "250m"
      memory = "256Mi"
    }
    requests = {
      cpu    = "100m"
      memory = "128Mi"
    }
  }
}

