variable "openai_api_key" {
  type        = string
  description = "The OpenAI API key"
}

variable "email_sender" {
  type        = string
  description = "The email of where the tales are sent from"
}

variable "email_destination" {
  type        = string
  description = "The email of where to send your tales"
}
