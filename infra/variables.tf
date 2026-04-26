variable "aws_region" {
  type    = string
  default = "eu-west-3"
}

variable "cloudfront_origin" {
  type        = string
  description = "Arcana SPA CloudFront URL allowed in CORS (e.g. https://xxx.cloudfront.net)"
}

variable "scw_api_url" {
  type        = string
  description = "Scaleway inference endpoint URL"
  sensitive   = true
}

variable "scw_secret_key" {
  type        = string
  description = "Scaleway secret API key"
  sensitive   = true
}

variable "dynamodb_endpoint" {
  type        = string
  description = "DynamoDB endpoint override — set to http://dynamodb-local:8000 for local dev, leave empty for AWS"
  default     = ""
}

variable "tags" {
  description = "Tags applied to all AWS ressources"
  type        = map(string)
  default = {
    Project     = "arcana"
    Module      = "ai"
    Environment = "production"
    ManagedBy   = "terraform"
  }
}
