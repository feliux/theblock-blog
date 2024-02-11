resource "aws_kms_key" "kms" {
  description              = var.kms.description
  key_usage                = "ENCRYPT_DECRYPT"
  customer_master_key_spec = "SYMMETRIC_DEFAULT"
  deletion_window_in_days  = 7
  enable_key_rotation      = false
  tags                     = var.custom_tags
}

resource "aws_kms_alias" "kms_alias" {
  name          = "alias/${var.kms.name}"
  target_key_id = aws_kms_key.kms.key_id
}

resource "aws_kms_grant" "lambda" {
  name              = var.kms.grant_name
  key_id            = aws_kms_key.kms.key_id
  grantee_principal = aws_iam_role.lambda_role.arn
  operations        = ["Encrypt", "Decrypt", "GenerateDataKey"] // GenerateDataKey for saving on s3
}

resource "aws_secretsmanager_secret" "sm" {
  name        = var.lambda_conf.lambda_environment.SECRET_NAME // var.sm.name
  description = var.sm.description
  kms_key_id  = aws_kms_key.kms.id
}

resource "aws_secretsmanager_secret_version" "secrets" {
  secret_id     = aws_secretsmanager_secret.sm.id
  secret_string = jsonencode(var.secrets)
}
