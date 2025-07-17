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

variable "cfm-agent_image" {
  description = "Docker image"
  type        = string
}

variable "pull_policy" {
  description = "Docker image pull policy"
  type        = string
}

variable "log_level" {
  description = "Log level"
  type        = string
  default     = "info"
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

variable "cfm-agent_port" {
  description = "Port that cfm-agent HTTP server listens on"
  type        = number
  default     = 8080
}

variable "metrics_port" {
  description = "Port that cfm-agent metrics server listens on"
  type        = number
  default     = 9090
}

variable "enable_nodeport" {
  description = "Enable NodePort service for external access"
  type        = bool
  default     = false
}

variable "cfm-agent_nodeport" {
  description = "NodePort HTTP server external access"
  type        = number
  default     = 30080
}

variable "metrics_nodeport" {
  description = "NodePort metrics server external access"
  type        = number
  default     = 30090
}

variable "fulcrum_core_service" {
  description = "Service name for Fulcrum core"
  type = string
}

variable "fulcrum_core_port" {
  description = "Port for Fulcrum core"
  type = number
}

variable "pmanager_service_url" {
  description = "URL for CFM Provision Manager"
  type = string
}

variable "tmanager_service_url" {
  description = "URL for CFM Tenant Manager"
  type = string
}

