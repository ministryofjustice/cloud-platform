resource "aws_s3_bucket" "example" {
  bucket = "my-test-bucket"
  tags = {
    Name        = "example"
    Environment = "test"
    Owner       = "team"
  }
}

variable "region" {
  description = "AWS region"
  default     = "eu-west-2"
  type        = string
}

output "bucket_id" {
  value = aws_s3_bucket.example.id
}
