resource "aws_s3_bucket" "example" {
bucket = "my-bucket"
  tags = {
    Name        = "my-bucket"
  Environment = "dev"
  }
}
