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

locals {
  default_labels = {
    app = "cfm-agent"
  }
  labels = merge(local.default_labels, var.labels)
}

variable "fulcrum_token" {
  description = "Fulcrum API token"
  type        = string
}

resource "kubernetes_deployment" "cfm-agent" {
  metadata {
    name      = "cfm-agent-server"
    namespace = var.namespace
    labels    = local.labels
  }

  spec {
    replicas = var.replicas

    selector {
      match_labels = local.default_labels
    }

    template {
      metadata {
        labels = local.labels
      }

      spec {
        container {
          image             = var.cfm-agent_image
          name              = "cfm-agent"
          image_pull_policy = var.pull_policy

          # port {
          #   container_port = var.cfm-agent_port
          #   name           = "http"
          # }
          #
          # port {
          #   container_port = var.metrics_port
          #   name           = "metrics"
          # }

          env {
            name  = "CFM-AGENT_TMANAGER_URL"
            value = var.tmanager_service_url
          }

          env {
            name  = "CFM-AGENT_PMANAGER_URL"
            value = var.pmanager_service_url
          }

          env {
            name  = "CFM-AGENT_FULCRUM_URI"
            value = "http://${var.fulcrum_core_service}:${var.fulcrum_core_port}"
          }


          env {
            name  = "CFM-AGENT_FULCRUM_TOKEN"
            value = var.fulcrum_token
          }

          resources {
            limits   = var.resources.limits
            requests = var.resources.requests
          }

          # liveness_probe {
          #   http_get {
          #     path = "/health"
          #     port = var.cfm-agent_port
          #   }
          #   initial_delay_seconds = 30
          #   period_seconds        = 10
          #   timeout_seconds       = 5
          #   failure_threshold     = 3
          # }
          #
          # readiness_probe {
          #   http_get {
          #     path = "/ready"
          #     port = var.cfm-agent_port
          #   }
          #   initial_delay_seconds = 5
          #   period_seconds        = 5
          #   timeout_seconds       = 3
          #   failure_threshold     = 3
          # }
          #
          # startup_probe {
          #   http_get {
          #     path = "/health"
          #     port = var.cfm-agent_port
          #   }
          #   initial_delay_seconds = 10
          #   period_seconds        = 10
          #   timeout_seconds       = 3
          #   failure_threshold     = 10
          # }
        }
      }
    }
  }
}

# cfm-agent Service
resource "kubernetes_service" "cfm-agent" {
  metadata {
    name      = "cfm-agent-service"
    namespace = var.namespace
    labels    = local.labels
  }

  spec {
    selector = local.default_labels

    port {
      name        = "http"
      port        = var.cfm-agent_port
      target_port = var.cfm-agent_port
    }

    port {
      name        = "metrics"
      port        = var.metrics_port
      target_port = var.metrics_port
    }

    type = "ClusterIP"
  }
}

# NodePort service for external access
resource "kubernetes_service" "cfm-agent_nodeport" {
  count = var.enable_nodeport ? 1 : 0

  metadata {
    name      = "cfm-agent-nodeport"
    namespace = var.namespace
    labels    = local.labels
  }

  spec {
    selector = local.default_labels

    port {
      name        = "http"
      port        = var.cfm-agent_port
      target_port = var.cfm-agent_port
      node_port   = var.cfm-agent_nodeport
    }

    port {
      name        = "metrics"
      port        = var.metrics_port
      target_port = var.metrics_port
      node_port   = var.metrics_nodeport
    }

    type = "NodePort"
  }
}