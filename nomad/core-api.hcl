job "awanku-stack-core-api" {
    datacenters = ["dc1"]
    group "core-api" {
        ephemeral_disk {
            migrate = true
            sticky  = true
            size    = "500"
        }
        task "core-api" {
            driver = "docker"
            config {
                image = "docker.awanku.id/awanku/core-api:latest"
                force_pull = true
                auth {
                    username = "awanku"
                    password = "rahasia"
                }
                port_map {
                    http = 3000
                }
            }
            service {
                name = "awanku-core-api"
                port = "http"
                check {
                    type = "http"
                    path = "/status"
                    port = "http"
                    interval = "10s"
                    timeout = "1s"
                    check_restart {
                        limit = 3
                        grace = "30s"
                    }
                }
                tags = [
                    "traefik.enable=true",
                    "traefik.http.routers.awanku-stack-core-api-http.rule=Host(`api.awanku.id`)",
                    "traefik.http.routers.awanku-stack-core-api-http.entrypoints=http",
                    "traefik.http.routers.awanku-stack-core-api-http.middlewares=httpToHttps@consul",
                    "traefik.http.routers.awanku-stack-core-api-https.rule=Host(`api.awanku.id`)",
                    "traefik.http.routers.awanku-stack-core-api-https.entrypoints=https",
                    "traefik.http.routers.awanku-stack-core-api-https.tls=true",
                    "traefik.http.routers.awanku-stack-core-api-https.tls.options=default",
                ]
            }
            env {
                ENVIRONMENT = "production"
                DATABASE_URL = "postgres://awanku:rahasia@${NOMAD_IP_maindb_pg}:${NOMAD_PORT_maindb_pg}/awanku?sslmode=disable"
                OAUTH_SECRET_KEY = "supersecretkey"
                GITHUB_APP_ID = 73537
                GITHUB_APP_PRIVATE_KEY_PATH = "/local/github.private-key.pem"
                GITHUB_APP_INSTALL_URL = "https://github.com/apps/awanku-development/installations/new"
            }
            template {
                data = "{{ key \"awanku/credentials/github/dev.private-key.pem\" }}"
                destination = "local/github.private-key.pem"
            }
            resources {
                network {
                    port "http" {}
                }
            }
            meta {
                VERSION = "current_version"
            }
        }
        task "maindb" {
            driver = "docker"
            config {
                image = "postgres:12"
                port_map {
                    pg = 5432
                }
            }
            service {
                name = "awanku-maindb"
                port = "pg"
                check {
                    type     = "tcp"
                    port     = "pg"
                    interval = "10s"
                    timeout  = "1s"
                    check_restart {
                        limit = 3
                        grace = "30s"
                    }
                }
            }
            env {
                POSTGRES_USER = "awanku"
                POSTGRES_PASSWORD = "rahasia"
                POSTGRES_DB = "awanku"
                PGDATA = "/alloc/data/postgres/pgdata"
            }
            resources {
                network {
                    port "pg" {}
                }
            }
        }
    }
}
