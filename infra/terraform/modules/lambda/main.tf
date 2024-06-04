data "aws_iam_policy_document" "assume_lambda_role" {
  statement {
    actions = ["sts:AssumeRole"]

    principals {
      type        = "Service"
      identifiers = ["lambda.amazonaws.com"]
    }
  }
}

resource "null_resource" "app_binary" {
  triggers = {
    always_run = var.recompile_go
  }

  provisioner "local-exec" {
    command = "GOOS=linux GOARCH=arm64 CGO_ENABLED=0 GOFLAGS=-trimpath go build -mod=readonly -ldflags='-s -w' -o ${var.app_binary_path} ${var.app_src_path}"
  }
}

data "archive_file" "app_archive" {
  depends_on = [null_resource.app_binary]

  type        = "zip"
  source_file = var.app_binary_path
  output_path = var.app_archive_path
}

resource "aws_lambda_function" "app" {
  function_name = var.function_name
  description   = var.description
  role          = aws_iam_role.lambda.arn
  handler       = "bootstrap"
  runtime       = "provided.al2023"
  architectures = ["arm64"]
  memory_size   = var.memory_size
  timeout       = var.timeout

  ephemeral_storage {
    size = 512
  }

  filename         = var.app_archive_path
  source_code_hash = data.archive_file.app_archive.output_base64sha256

  environment {
    variables = var.environment
  }
}

resource "aws_iam_role" "lambda" {
  name               = "${var.function_name}-AssumeLambdaRole-${terraform.workspace}"
  description        = "Role for lambda to assume its execution role. Grant the Lambda service principal permission to assume our role"
  assume_role_policy = data.aws_iam_policy_document.assume_lambda_role.json
}

resource "aws_iam_policy" "lambda" {
  name        = "${var.function_name}-lambda-permissions-${terraform.workspace}"
  path        = "/"
  description = "For Lambda and what it can access"

  #  policy = data.aws_iam_policy_document.lambda.json

  policy = jsonencode({
    Version = "2012-10-17"
    Statement = concat([
      {
        Action = [
          "logs:CreateLogGroup",
          "logs:CreateLogStream",
          "logs:PutLogEvents",
        ]
        Effect   = "Allow"
        Resource = "arn:aws:logs:*:*:*"
      }
    ],
      [for perm_name, perm in var.permissions : {
        Action   = perm.actions
        Effect   = perm.effect
        Resource = perm.resources
      }]
    )
  })
}

resource "aws_iam_role_policy_attachment" "lambda" {
  role       = aws_iam_role.lambda.name
  policy_arn = aws_iam_policy.lambda.arn
}