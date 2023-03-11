locals {
  project      = "tail-time"
  src_path     = "../../cmd/tail-time-lambda"
  binary_path  = "./dist/bin/${local.project}"
  binary_name  = local.project
  archive_path = "./dist/app.zip"
}

// S3

resource "aws_s3_bucket" "tales" {
  bucket = "${local.project}-tales"
}

resource "aws_s3_bucket_acl" "tales" {
  bucket = aws_s3_bucket.tales.id
  acl    = "private"
}

// Lambda

resource "null_resource" "app_binary" {
#  triggers = {
#    always_run = timestamp() // TODO fix this
#  }

  provisioner "local-exec" {
    command = "GOOS=linux GOARCH=amd64 CGO_ENABLED=0 GOFLAGS=-trimpath go build -mod=readonly -ldflags='-s -w' -o ${local.binary_path} ${local.src_path}"
  }
}

data "archive_file" "app_archive" {
  depends_on = [null_resource.app_binary]

  type        = "zip"
  source_file = local.binary_path
  output_path = local.archive_path
}

resource "aws_lambda_function" "app" {
  function_name = local.project
  description   = "Haha"
  role          = aws_iam_role.lambda.arn
  handler       = local.binary_name
  runtime       = "go1.x"
  memory_size   = 128
  timeout       = 120

  filename         = local.archive_path
  source_code_hash = data.archive_file.app_archive.output_base64sha256

  environment {
    variables = {
      OPENAI_API_KEY = ""
    }
  }
}

#resource "aws_lambda_function_url" "app" {
#  function_name      = aws_lambda_function.app.function_name
#  authorization_type = "NONE"
#
#  cors {
#    allow_credentials = false
#    allow_origins     = ["*"]
#    allow_methods     = ["GET"]
#    allow_headers     = ["date", "keep-alive"]
#    expose_headers    = ["keep-alive", "date"]
#    max_age           = 86400
#  }
#}

data "aws_iam_policy_document" "assume_lambda_role" {
  statement {
    actions = ["sts:AssumeRole"]

    principals {
      type        = "Service"
      identifiers = ["lambda.amazonaws.com"]
    }
  }
}

resource "aws_iam_role" "lambda" {
  name               = "AssumeLambdaRole"
  description        = "Role for lambda to assume its execution role. Grant the Lambda service principal permission to assume our role"
  assume_role_policy = data.aws_iam_policy_document.assume_lambda_role.json
}

resource "aws_iam_policy" "lambda" {
  name        = "Lambda-permissions"
  path        = "/"
  description = "For Lambda and what it can access"

  policy = jsonencode({
    Version = "2012-10-17"
    Statement = [
      {
        Action = [
          "s3:PutObject",
        ]
        Effect   = "Allow"
        Resource = "${aws_s3_bucket.tales.arn}/*"
      },
    ]
  })
}

resource aws_iam_role_policy_attachment lambda {
  role       = aws_iam_role.lambda.name
  policy_arn = aws_iam_policy.lambda.arn
}

