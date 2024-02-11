resource "aws_s3_bucket" "s3" {
  bucket = var.lambda_conf.lambda_environment.BUCKET_NAME
  // acl           = "private" // by default - deprecated
  force_destroy = true
  tags          = var.custom_tags
}

resource "aws_s3_object" "avatars" {
  bucket        = aws_s3_bucket.s3.id
  key           = "avatars/"
  force_destroy = true
}

resource "aws_s3_object" "banners" {
  bucket        = aws_s3_bucket.s3.id
  key           = "banners/"
  force_destroy = true
}
