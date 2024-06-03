locals {
  project               = var.project_name
  generate_src_path     = "../../cmd/${local.project}-lambda-generate"
  generate_binary_path  = "./dist/bin/generate/bootstrap"
  generate_archive_path = "./dist/generate.zip"

  send_src_path     = "../../cmd/${local.project}-lambda-send"
  send_binary_path  = "./dist/bin/send/bootstrap"
  send_archive_path = "./dist/send.zip"

  generate_audio_src_path     = "../../cmd/${local.project}-lambda-generate-audio"
  generate_audio_binary_path  = "./dist/bin/generate-audio/bootstrap"
  generate_audio_archive_path = "./dist/generate-audio.zip"

  #   alwaysDeployGo = timestamp()
  alwaysDeployGo = 1
}

// S3

resource "aws_s3_bucket" "tales" {
  bucket = "${local.project}-tales-${terraform.workspace}-73d2d65dca41"
}

// SES

resource "aws_ses_email_identity" "email" {
  email = var.email_from
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
    always_run = local.alwaysDeployGo
  }

  provisioner "local-exec" {
    command = "GOOS=linux GOARCH=arm64 CGO_ENABLED=0 GOFLAGS=-trimpath go build -mod=readonly -ldflags='-s -w' -o ${local.generate_binary_path} ${local.generate_src_path}"
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
  handler       = "bootstrap"
  architectures = ["arm64"]
  runtime       = "provided.al2023"
  memory_size   = 128
  timeout       = 120

  filename         = local.generate_archive_path
  source_code_hash = data.archive_file.generate_app_archive.output_base64sha256

  environment {
    variables = {
      OPENAI_API_KEY            = var.openai_api_key
      DESTINATION_BUCKET_NAME   = aws_s3_bucket.tales.id
      DESTINATION_BUCKET_REGION = aws_s3_bucket.tales.region
    }
  }
}

resource "aws_iam_role" "generate_lambda" {
  name               = "AssumeGenerateLambdaRole-${terraform.workspace}"
  description        = "Role for lambda to assume its execution role. Grant the Lambda service principal permission to assume our role"
  assume_role_policy = data.aws_iam_policy_document.assume_lambda_role.json
}

resource "aws_iam_policy" "generate_lambda" {
  name        = "Generate-lambda-permissions-${terraform.workspace}"
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

// Lambda to generate tale audio
// TODO make Lambda module

resource "null_resource" "generate_audio_app_binary" {
  triggers = {
    always_run = local.alwaysDeployGo
  }

  provisioner "local-exec" {
    command = "GOOS=linux GOARCH=arm64 CGO_ENABLED=0 GOFLAGS=-trimpath go build -mod=readonly -ldflags='-s -w' -o ${local.generate_audio_binary_path} ${local.generate_audio_src_path}"
  }
}

data "archive_file" "generate_audio_app_archive" {
  depends_on = [null_resource.generate_audio_app_binary]

  type        = "zip"
  source_file = local.generate_audio_binary_path
  output_path = local.generate_audio_archive_path
}

resource "aws_lambda_function" "generate_audio_app" {
  function_name = "${local.project}-generate-audio"
  description   = "Lambda to generate tale audio and store in s3"
  role          = aws_iam_role.generate_audio_lambda.arn
  handler       = "bootstrap"
  architectures = ["arm64"]
  runtime       = "provided.al2023"
  memory_size   = 128
  timeout       = 240

  filename         = local.generate_audio_archive_path
  source_code_hash = data.archive_file.generate_audio_app_archive.output_base64sha256

  environment {
    variables = {
      OPENAI_API_KEY            = var.openai_api_key
      DESTINATION_BUCKET_NAME   = aws_s3_bucket.tales.id
      DESTINATION_BUCKET_REGION = aws_s3_bucket.tales.region
    }
  }
}

resource "aws_iam_role" "generate_audio_lambda" {
  name               = "AssumeGenerateAudioLambdaRole-${terraform.workspace}"
  description        = "Role for lambda to assume its execution role. Grant the Lambda service principal permission to assume our role"
  assume_role_policy = data.aws_iam_policy_document.assume_lambda_role.json
}

resource "aws_iam_policy" "generate_audio_lambda" {
  name        = "Generate-audio-lambda-permissions-${terraform.workspace}"
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
          "s3:GetObject",
        ]
        Effect   = "Allow"
        Resource = "${aws_s3_bucket.tales.arn}/*"
      }
    ]
  })
}

resource "aws_iam_role_policy_attachment" "generate_audio_lambda" {
  role       = aws_iam_role.generate_audio_lambda.name
  policy_arn = aws_iam_policy.generate_audio_lambda.arn
}


// Scheduled trigger

resource "aws_cloudwatch_event_rule" "every_day" {
  name                = "every-day"
  description         = "Fires once a day"
  schedule_expression = "cron(0 14 * * ? *)"
}

resource "aws_cloudwatch_event_target" "every_day" {
  rule  = aws_cloudwatch_event_rule.every_day.name
  arn   = aws_lambda_function.generate_app.arn
  input = jsonencode({ topic : var.default_tale_topic })
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
    always_run = local.alwaysDeployGo
  }

  provisioner "local-exec" {
    command = "GOOS=linux GOARCH=arm64 CGO_ENABLED=0 GOFLAGS=-trimpath go build -mod=readonly -ldflags='-s -w' -o ${local.send_binary_path} ${local.send_src_path}"
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
  handler       = "bootstrap"
  runtime       = "provided.al2023"
  architectures = ["arm64"]
  memory_size   = 128
  timeout       = 120

  ephemeral_storage {
    size = 512
  }

  filename         = local.send_archive_path
  source_code_hash = data.archive_file.send_app_archive.output_base64sha256

  environment {
    variables = {
      EMAIL_FROM           = var.email_from
      EMAIL_TO             = var.email_to
      SOURCE_BUCKET_REGION = aws_s3_bucket.tales.region
    }
  }
}

resource "aws_iam_role" "send_lambda" {
  name               = "AssumeSendLambdaRole-${terraform.workspace}"
  description        = "Role for lambda to assume its execution role. Grant the Lambda service principal permission to assume our role"
  assume_role_policy = data.aws_iam_policy_document.assume_lambda_role.json
}

resource "aws_iam_policy" "send_lambda" {
  name        = "Send-lambda-permissions-${terraform.workspace}"
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
          "ses:SendRawEmail",
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

// Routing

resource "aws_cloudwatch_event_rule" "tale_text_created_event_rule" {
  name = "text-created-${terraform.workspace}"
  event_pattern = jsonencode({
    source : ["aws.s3"],
    detail-type : ["Object Created"],
    "detail" : {
      "bucket" : {
        "name" : [aws_s3_bucket.tales.id]
      },
      #       "object" : {
      #         "key" : [{ "prefix" : "raw/" }]
      #       }
    }
  })
}

resource "aws_cloudwatch_event_target" "send_text_email" {
  rule      = aws_cloudwatch_event_rule.tale_text_created_event_rule.name
  target_id = "send-text-email"
  arn       = aws_lambda_function.send_app.arn
}

resource "aws_s3_bucket_notification" "notify_event_bus" {
  bucket      = aws_s3_bucket.tales.id
  eventbridge = true
}
