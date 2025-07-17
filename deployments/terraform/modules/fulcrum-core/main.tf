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

resource "kubernetes_deployment" "fulcrum_core" {
  metadata {
    name      = "fulcrum-core"
    namespace = var.namespace
    labels = merge(var.labels, {
      app = "fulcrum-core"
    })
  }

  spec {
    replicas = var.replicas

    selector {
      match_labels = {
        app = "fulcrum-core"
      }
    }

    template {
      metadata {
        labels = merge(var.labels, {
          app = "fulcrum-core"
        })
      }

      spec {
        container {
          image             = var.fulcrum_core_image
          name              = var.fulcrum_core_service
          image_pull_policy = var.pull_policy

          port {
            container_port = var.fulcrum_core_port
          }

          env {
            name  = "FULCRUM_DB_DSN"
            value = local.fulcrum_db_dsn
          }

          env {
            name  = "FULCRUM_DB_LOG_LEVEL"
            value = var.fulcrum_db_log_level
          }

          env {
            name  = "FULCRUM_DB_LOG_FORMAT"
            value = var.fulcrum_db_log_format
          }

          resources {
            limits = {
              cpu    = var.resources.limits.cpu
              memory = var.resources.limits.memory
            }
            requests = {
              cpu    = var.resources.requests.cpu
              memory = var.resources.requests.memory
            }
          }

          liveness_probe {
            http_get {
              path = "/health"
              port = 8081
            }
            initial_delay_seconds = 3000
            period_seconds        = 1000
          }

          readiness_probe {
            http_get {
              path = "/ready"
              port = 8081
            }
            initial_delay_seconds = 5
            period_seconds        = 5
          }
        }
      }
    }
  }
}

resource "kubernetes_service" "fulcrum_core" {
  metadata {
    name      = "fulcrum-core"
    namespace = var.namespace
    labels = merge(var.labels, {
      app = "fulcrum-core"
    })
  }

  spec {
    selector = {
      app = "fulcrum-core"
    }

    port {
      port        = var.fulcrum_core_port
      target_port = var.fulcrum_core_port
      node_port   = var.fulcrum_core_nodeport
    }

    type = "NodePort"
  }
}