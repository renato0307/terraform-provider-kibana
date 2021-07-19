terraform {
  required_providers {
    kibana = {
      source  = "renatoalvestorres.net/terraform/kibana"
      version = ">= 0.0.2"
    }
  }
}

provider "kibana" {
  host     = var.kibana_url
  username = var.kibana_username
  password = var.kibana_password
  space    = var.kibana_space
}
