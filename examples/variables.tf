


variable "kibana_url" {
  description = "Endpoint to access Kibana"
  type        = string
}

variable "kibana_username" {
  description = "Username to access Kibana"
  sensitive   = true
  type        = string
}

variable "kibana_password" {
  description = "Password to access Kibana"
  sensitive   = true
  type        = string
}

variable "kibana_space" {
  description = "Kibana space to use"
  type        = string
}
