job "awanku-stack-api-docs" {
    datacenters = ["dc1"]
    group "api-docs" {
        task "api-docs" {
            driver = "docker"
            config {
                image = "docker.awanku.id/awanku/core-api-docs:latest"
                force_pull = true
                auth {
                    username = "awanku"
                    password = "rahasia"
                }
                port_map {
                    http = 80
                }
            }
            service {
                name = "awanku-core-api-docs"
                port = "http"
                check {
                    type     = "http"
                    port     = "http"
                    path     = "/docs/"
                    interval = "10s"
                    timeout  = "1s"
                    check_restart {
                        limit = 3
                        grace = "30s"
                    }
                }
                tags = [
                    "traefik.enable=true",
                    "traefik.http.routers.awanku-stack-core-api-docs-https.rule=Host(`api.awanku.id`) && PathPrefix(`/docs/`)",
                    "traefik.http.routers.awanku-stack-core-api-docs-https.entrypoints=https",
                    "traefik.http.routers.awanku-stack-core-api-docs-https.tls=true",
                    "traefik.http.routers.awanku-stack-core-api-docs-https.tls.options=default",
                ]
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

    }
}
