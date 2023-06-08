variable "project_name" {
  type        = string
  description = "The name of the project, used for prefixes, naming etc"
}

variable "aws_region" {
  type        = string
  description = "The default region of AWS resources"
}

variable "openai_api_key" {
  type        = string
  description = "The OpenAI API key"
}

variable "email_from" {
  type        = string
  description = "The email of where the tales are sent from"
}

variable "email_to" {
  type        = string
  description = "The email of where to send your tales"
}

variable "default_tale_topic" {
  type        = string
  description = "Default topic of tales"
}
