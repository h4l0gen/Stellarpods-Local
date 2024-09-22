# Configure the Kubernetes provider
provider "kubernetes" {
  host = "https://192.168.58.2:8443"  # Use the IP from your kubectl cluster-info output

  client_certificate     = file("~/.minikube/profiles/minikube/client.crt")
  client_key             = file("~/.minikube/profiles/minikube/client.key")
  cluster_ca_certificate = file("~/.minikube/ca.crt")
}

# Create a namespace for our CTF challenges
resource "kubernetes_namespace" "ctf_challenges" {
  metadata {
    name = "ctf-challenges"
  }
}

# Deploy a CTF challenge (example)
resource "kubernetes_pod" "example_challenge" {
  metadata {
    name      = "example-ctf-challenge"
    namespace = kubernetes_namespace.ctf_challenges.metadata[0].name
  }

  spec {
    container {
      image = "nginx:1.19.4"  # Replace with your actual challenge image
      name  = "ctf-challenge-container"

      port {
        container_port = 80
      }
    }
  }

  depends_on = [kubernetes_namespace.ctf_challenges]
}

# Expose the challenge with a service
resource "kubernetes_service" "example_challenge_service" {
  metadata {
    name      = "example-ctf-challenge-service"
    namespace = kubernetes_namespace.ctf_challenges.metadata[0].name
  }

  spec {
    selector = {
      app = kubernetes_pod.example_challenge.metadata[0].name
    }

    port {
      port        = 80
      target_port = 80
    }

    type = "NodePort"
  }

  depends_on = [kubernetes_pod.example_challenge]
}

# Output the URL to access the challenge
output "challenge_url" {
  value = "http://${data.external.minikube_ip.result.ip}:${kubernetes_service.example_challenge_service.spec[0].port[0].node_port}"
}

# Data source to get Minikube IP
data "external" "minikube_ip" {
  program = ["sh", "-c", "echo '{\"ip\": \"'$(minikube ip)'\"}'"]
}