variable "function_name" {
  type = string
}

variable "description" {
  type = string
}

variable "app_binary_path" {
  type = string
}

variable "app_archive_path" {
  type = string
}

variable "app_src_path" {
  type = string
}

variable "environment" {
  type        = map(string)
  default     = {}
}

variable "permissions" {
  type        = map(object({
    actions   = list(string)
    effect   = string
    resources = list(string)
  }))
  default     = {}
}

variable "memory_size" {
  type = number
  default = 128
}

variable "timeout" {
  type = number
  default = 60
}