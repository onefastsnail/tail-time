locals {
  project               = "tail-time"
  generate_src_path     = "../../cmd/${local.project}-lambda-generate"
  generate_binary_path  = "./dist/bin/${local.project}-generate"
  generate_binary_name  = "${local.project}-generate"
  generate_archive_path = "./dist/generate.zip"

  send_src_path     = "../../cmd/${local.project}-lambda-send"
  send_binary_path  = "./dist/bin/${local.project}-send"
  send_binary_name  = "${local.project}-send"
  send_archive_path = "./dist/send.zip"
}

// S3

resource "aws_s3_bucket" "tales" {
  bucket = "${local.project}-tales-${terraform.workspace}"
}

resource "aws_s3_bucket_acl" "tales" {
  bucket = aws_s3_bucket.tales.id
  acl    = "private"
}

// SES

resource "aws_ses_email_identity" "email" {
  email = var.email_sender
}

// IAM

data "aws_iam_policy_document" "assume_lambda_role" {
  statement {
    actions = ["sts:AssumeRole"]

    principals {
      type        = "Service"
      identifiers = ["lambda.amazonaws.com"]
    }
  }
}

// Lambda to generate tales

resource "null_resource" "generate_app_binary" {
  triggers = {
    always_run = timestamp() // TODO fix this
  }

  provisioner "local-exec" {
    command = "GOOS=linux GOARCH=amd64 CGO_ENABLED=0 GOFLAGS=-trimpath go build -mod=readonly -ldflags='-s -w' -o ${local.generate_binary_path} ${local.generate_src_path}"
  }
}

data "archive_file" "generate_app_archive" {
  depends_on = [null_resource.generate_app_binary]

  type        = "zip"
  source_file = local.generate_binary_path
  output_path = local.generate_archive_path
}

resource "aws_lambda_function" "generate_app" {
  function_name = "${local.project}-generate"
  description   = "Lambda to generate tales and store in s3"
  role          = aws_iam_role.generate_lambda.arn
  handler       = local.generate_binary_name
  runtime       = "go1.x"
  memory_size   = 128
  timeout       = 120

  filename         = local.generate_archive_path
  source_code_hash = data.archive_file.generate_app_archive.output_base64sha256

  environment {
    variables = {
      OPENAI_API_KEY          = var.openai_api_key
      DESTINATION_BUCKET_NAME = aws_s3_bucket.tales.id
    }
  }
}

resource "aws_iam_role" "generate_lambda" {
  name               = "AssumeGenerateLambdaRole"
  description        = "Role for lambda to assume its execution role. Grant the Lambda service principal permission to assume our role"
  assume_role_policy = data.aws_iam_policy_document.assume_lambda_role.json
}

resource "aws_iam_policy" "generate_lambda" {
  name        = "Generate-lambda-permissions"
  path        = "/"
  description = "For Lambda and what it can access"

  policy = jsonencode({
    Version = "2012-10-17"
    Statement = [
      {
        Action = [
          "logs:CreateLogGroup",
          "logs:CreateLogStream",
          "logs:PutLogEvents",
        ]
        Effect   = "Allow"
        Resource = "arn:aws:logs:*:*:*"
      },
      {
        Action = [
          "s3:PutObject",
        ]
        Effect   = "Allow"
        Resource = "${aws_s3_bucket.tales.arn}/*"
      }
    ]
  })
}

resource "aws_iam_role_policy_attachment" "generate_lambda" {
  role       = aws_iam_role.generate_lambda.name
  policy_arn = aws_iam_policy.generate_lambda.arn
}


// Scheduled trigger

resource "aws_cloudwatch_event_rule" "every_day" {
  name                = "every-day"
  description         = "Fires every once a day"
  schedule_expression = "cron(0 17 * * ? *)"
}

resource "aws_cloudwatch_event_target" "every_day" {
  rule = aws_cloudwatch_event_rule.every_day.name
  arn  = aws_lambda_function.generate_app.arn
}

resource "aws_lambda_permission" "allow_cloudwatch_to_call_generate" {
  statement_id  = "AllowExecutionFromCloudWatch"
  action        = "lambda:InvokeFunction"
  function_name = aws_lambda_function.generate_app.function_name
  principal     = "events.amazonaws.com"
  source_arn    = aws_cloudwatch_event_rule.every_day.arn
}

// Lambda to send out tales

resource "null_resource" "send_app_binary" {
  triggers = {
    always_run = timestamp() // TODO fix this
  }

  provisioner "local-exec" {
    command = "GOOS=linux GOARCH=amd64 CGO_ENABLED=0 GOFLAGS=-trimpath go build -mod=readonly -ldflags='-s -w' -o ${local.send_binary_path} ${local.send_src_path}"
  }
}

data "archive_file" "send_app_archive" {
  depends_on = [null_resource.send_app_binary]

  type        = "zip"
  source_file = local.send_binary_path
  output_path = local.send_archive_path
}

resource "aws_lambda_function" "send_app" {
  function_name = "${local.project}-send"
  description   = "Lambda to send tales from tales being stored in s3"
  role          = aws_iam_role.send_lambda.arn
  handler       = local.send_binary_name
  runtime       = "go1.x"
  memory_size   = 128
  timeout       = 120

  filename         = local.send_archive_path
  source_code_hash = data.archive_file.send_app_archive.output_base64sha256

  environment {
    variables = {
      EMAIL_DESTINATION = var.email_destination
    }
  }
}

resource "aws_iam_role" "send_lambda" {
  name               = "AssumeSendLambdaRole"
  description        = "Role for lambda to assume its execution role. Grant the Lambda service principal permission to assume our role"
  assume_role_policy = data.aws_iam_policy_document.assume_lambda_role.json
}

resource "aws_iam_policy" "send_lambda" {
  name        = "Send-lambda-permissions"
  path        = "/"
  description = "For Lambda and what it can access"

  #  policy = data.aws_iam_policy_document.send_lambda.json

  policy = jsonencode({
    Version = "2012-10-17"
    Statement = [
      {
        Action = [
          "logs:CreateLogGroup",
          "logs:CreateLogStream",
          "logs:PutLogEvents",
        ]
        Effect   = "Allow"
        Resource = "arn:aws:logs:*:*:*"
      },
      {
        Action = [
          "s3:GetObject",
        ]
        Effect   = "Allow"
        Resource = "${aws_s3_bucket.tales.arn}/*"
      },
      {
        Action = [
          "ses:SendEmail",
        ]
        Effect   = "Allow"
        Resource = aws_ses_email_identity.email.arn
      },
    ]
  })
}

resource "aws_iam_role_policy_attachment" "send_lambda" {
  role       = aws_iam_role.send_lambda.name
  policy_arn = aws_iam_policy.send_lambda.arn
}

resource "aws_lambda_permission" "with_s3" {
  statement_id  = "AllowExecutionFromS3"
  action        = "lambda:InvokeFunction"
  function_name = aws_lambda_function.send_app.id
  principal     = "s3.amazonaws.com"
  source_arn    = aws_s3_bucket.tales.arn
}

resource "aws_s3_bucket_notification" "send" {
  bucket = aws_s3_bucket.tales.id

  lambda_function {
    lambda_function_arn = aws_lambda_function.send_app.arn
    events              = ["s3:ObjectCreated:*"]
    filter_prefix       = "raw/"
    filter_suffix       = ".txt"
  }
}
