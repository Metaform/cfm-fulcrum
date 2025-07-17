# PostgreSQL Deployment
resource "kubernetes_deployment" "postgres" {
  metadata {
    name      = "postgres"
    namespace = "default"
    labels = {
      app = "postgres"
    }
  }

  spec {
    replicas = var.replicas

    selector {
      match_labels = {
        app = "postgres"
      }
    }

    template {
      metadata {
        labels = {
          app = "postgres"
        }
      }

      spec {
        container {
          name  = "postgres"
          image = "postgres:17-alpine"

          env {
            name  = "POSTGRES_DB"
            value = "fulcrum_db"
          }

          env {
            name  = "POSTGRES_USER"
            value = var.postgres_user
          }

          env {
            name = "POSTGRES_PASSWORD"
            value_from {
              secret_key_ref {
                name = kubernetes_secret.postgres_secret.metadata[0].name
                key  = "password"
              }
            }
          }

          port {
            container_port = var.postgres_port
          }

          volume_mount {
            name       = "postgres-storage"
            mount_path = "/var/lib/postgresql/data"
          }

          volume_mount {
            name       = "postgres-init"
            mount_path = "/docker-entrypoint-initdb.d"
            read_only  = true
          }

          resources {
            limits   = var.resources.limits
            requests = var.resources.requests
          }

        }

        volume {
          name = "postgres-storage"
          empty_dir {}
        }

        volume {
          name = "postgres-init"
          config_map {
            name = kubernetes_config_map.postgres_init.metadata[0].name
          }
        }
      }
    }
  }

  depends_on = [kubernetes_secret.postgres_secret, kubernetes_config_map.postgres_init]
}

# PostgreSQL Service
resource "kubernetes_service" "postgres" {
  metadata {
    name      = var.postgres_service
    namespace = "default"
  }

  spec {
    selector = {
      app = "postgres"
    }

    port {
      port        = var.postgres_port
      target_port = 5432
    }

    type = "ClusterIP"
  }
}

# Postgres secret
resource "kubernetes_secret" "postgres_secret" {
  metadata {
    name      = "postgres-secret"
    namespace = "default"
  }

  data = {
    password = "password" # base64 encoded "password"
  }

  type = "Opaque"
}

# ConfigMap for PostgreSQL initialization
resource "kubernetes_config_map" "postgres_init" {
  metadata {
    name      = "postgres-init"
    namespace = "default"
  }
}