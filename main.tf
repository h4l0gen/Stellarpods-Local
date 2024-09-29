terraform {
  required_providers {
    kubernetes = {
      source  = "hashicorp/kubernetes"
      version = "~> 2.0"
    }
    docker = {
      source  = "kreuzwerker/docker"
      version = "~> 3.0.0"
    }
  }
}

provider "kubernetes" {
  host = "https://192.168.49.2:8443"  # Use the IP from your kubectl cluster-info output
  client_certificate     = file("~/.minikube/profiles/minikube/client.crt")
  client_key             = file("~/.minikube/profiles/minikube/client.key")
  cluster_ca_certificate = file("~/.minikube/ca.crt")
}

provider "docker" {}

variable "challenge_dir" {
  type    = string
  default = "challenge-github"  // here i did change
}

variable "download_path" {
  type    = string
  default = "."
}

resource "null_resource" "debug_info" {
  provisioner "local-exec" {
    command = <<-EOT
      echo "Debug Information:"
      echo "Download Path: ${var.download_path}"
      echo "Challenge Dir: ${var.challenge_dir}"
      echo "Full Path: ${var.download_path}/${var.challenge_dir}"
      echo "Files in Challenge Directory:"
      ls -la ${var.download_path}/${var.challenge_dir}
      echo "Dockerfile contents:"
      cat ${var.download_path}/${var.challenge_dir}/Dockerfile
    EOT
  }
}

resource "null_resource" "docker_build" {
  provisioner "local-exec" {
    command = "docker build -t ctf-challenge:latest ${var.download_path}/${var.challenge_dir}"
  }

  depends_on = [null_resource.debug_info]
}

resource "docker_image" "challenge_image" {
  name = "ctf-challenge:latest"
  keep_locally = false

  depends_on = [null_resource.docker_build]
}


resource "kubernetes_namespace" "ctf_challenges" {
  metadata {
    name = "ctf-challenges"
  }
}

# to load image to minikube from local docker daemon TODO changes
resource "null_resource" "load_image_to_minikube" {
  provisioner "local-exec" {
    command = "minikube image load ctf-challenge:latest"
  }
  depends_on = [docker_image.challenge_image]
}

resource "kubernetes_deployment" "challenge_deployment" {
  metadata {
    name      = "ctf-challenge"
    namespace = kubernetes_namespace.ctf_challenges.metadata[0].name
  }

  spec {
    replicas = 1

    selector {
      match_labels = {
        app = "ctf-challenge"
      }
    }

    template {
      metadata {
        labels = {
          app = "ctf-challenge"
        }
      }

      spec {
        container {
          image             = "ctf-challenge:latest"
          name              = "ctf-challenge-container"
          image_pull_policy = "Never"  # This tells Kubernetes to only use local images

          port {
            container_port = 80
          }
        }
      }
    }
  }

  depends_on = [kubernetes_namespace.ctf_challenges, docker_image.challenge_image]
}

resource "kubernetes_service" "challenge_service" {
  metadata {
    name      = "ctf-challenge-service"
    namespace = kubernetes_namespace.ctf_challenges.metadata[0].name
  }

  spec {
    selector = {
      app = kubernetes_deployment.challenge_deployment.spec[0].template[0].metadata[0].labels.app
    }

    port {
      port        = 80
      target_port = 80
    }

    type = "NodePort"
  }

  depends_on = [kubernetes_deployment.challenge_deployment]
}

data "external" "minikube_ip" {
  program = ["sh", "-c", "echo '{\"ip\": \"'$(minikube ip)'\"}'"]
}

output "challenge_url" {
  value = "http://${data.external.minikube_ip.result.ip}:${kubernetes_service.challenge_service.spec[0].port[0].node_port}"
}