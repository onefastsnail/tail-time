// S3

resource "aws_s3_bucket" "tales" {
  bucket = "${var.project_name}-tales-${terraform.workspace}-73d2d65dca41"
}

// SES

resource "aws_ses_email_identity" "email" {
  email = var.email_from
}

// Lambdas

module "generate_tale_text_lambda" {
  source           = "./modules/lambda"
  function_name    = "${var.project_name}-generate-tale-text"
  description      = "To generate tales from OpenAI"
  timeout          = 120
  app_src_path     = "../../cmd/${var.project_name}-lambda-generate-text"
  app_binary_path  = "./dist/bin/generate/bootstrap"
  app_archive_path = "./dist/generate.zip"
  permissions = {
    s3 = {
      actions   = ["s3:PutObject"]
      effect    = "Allow"
      resources = ["${aws_s3_bucket.tales.arn}/raw/*"]
    }
  }
  environment = {
    OPENAI_API_KEY            = var.openai_api_key
    DESTINATION_BUCKET_NAME   = aws_s3_bucket.tales.id
    DESTINATION_BUCKET_REGION = aws_s3_bucket.tales.region
  }
}

module "send_tale_text_lambda" {
  source           = "./modules/lambda"
  function_name    = "${var.project_name}-send-tale-text-email"
  description      = "To send tales as emails"
  timeout          = 120
  app_src_path     = "../../cmd/${var.project_name}-lambda-send-email"
  app_binary_path  = "./dist/bin/send/bootstrap"
  app_archive_path = "./dist/send.zip"
  permissions = {
    s3 = {
      actions   = ["s3:GetObject"]
      effect    = "Allow"
      resources = ["${aws_s3_bucket.tales.arn}/raw/*"]
    },
    ses = {
      actions   = ["ses:SendEmail", "ses:SendRawEmail"]
      effect    = "Allow"
      resources = [aws_ses_email_identity.email.arn]
    }
  }
  environment = {
    EMAIL_FROM           = var.email_from
    EMAIL_TO             = var.email_to
    SOURCE_BUCKET_REGION = aws_s3_bucket.tales.region
  }
}

module "generate_tale_audio_lambda" {
  source           = "./modules/lambda"
  function_name    = "${var.project_name}-generate-tale-audio"
  description      = "Generates audio versions of tales"
  timeout          = 120
  app_src_path     = "../../cmd/${var.project_name}-lambda-generate-audio"
  app_binary_path  = "./dist/bin/generate-audio/bootstrap"
  app_archive_path = "./dist/generate-audio.zip"
  permissions = {
    s3Put = {
      actions   = ["s3:PutObject"]
      effect    = "Allow"
      resources = ["${aws_s3_bucket.tales.arn}/audio/*"]
    }
    s3Get = {
      actions   = ["s3:GetObject"]
      effect    = "Allow"
      resources = ["${aws_s3_bucket.tales.arn}/raw/*"]
    }
  }
  environment = {
    OPENAI_API_KEY            = var.openai_api_key
    DESTINATION_BUCKET_NAME   = aws_s3_bucket.tales.id
    DESTINATION_BUCKET_REGION = aws_s3_bucket.tales.region
  }
}

// Routing

resource "aws_cloudwatch_event_rule" "every_day" {
  name                = "every-day"
  description         = "Fires once a day"
  schedule_expression = "cron(0 14 * * ? *)"
}

resource "aws_cloudwatch_event_target" "every_day" {
  rule  = aws_cloudwatch_event_rule.every_day.name
  arn   = module.generate_tale_text_lambda.lambda_function.arn
  input = jsonencode({ topic : var.default_tale_topic, language: "English" })
}

resource "aws_lambda_permission" "allow_cloudwatch_to_call_generate" {
  statement_id  = "AllowExecutionFromCloudWatch"
  action        = "lambda:InvokeFunction"
  function_name = module.generate_tale_text_lambda.lambda_function.function_name
  principal     = "events.amazonaws.com"
  source_arn    = aws_cloudwatch_event_rule.every_day.arn
}

resource "aws_lambda_permission" "with_s3" {
  statement_id  = "AllowExecutionFromS3"
  action        = "lambda:InvokeFunction"
  function_name = module.send_tale_text_lambda.lambda_function.id
  principal     = "s3.amazonaws.com"
  source_arn    = aws_s3_bucket.tales.arn
}

resource "aws_cloudwatch_event_rule" "tale_text_created_event_rule" {
  name = "text-created-${terraform.workspace}"
  event_pattern = jsonencode({
    source : ["aws.s3"],
    detail-type : ["Object Created"],
    "detail" : {
      "bucket" : {
        "name" : [aws_s3_bucket.tales.id]
      },
      "object" : {
        "key" : [{ "prefix" : "raw/" }]
      }
    }
  })
}

resource "aws_cloudwatch_event_target" "send_text_email" {
  rule      = aws_cloudwatch_event_rule.tale_text_created_event_rule.name
  target_id = "send-text-email"
  arn       = module.send_tale_text_lambda.lambda_function.arn
}

resource "aws_cloudwatch_event_target" "generate_audio_email" {
  rule      = aws_cloudwatch_event_rule.tale_text_created_event_rule.name
  target_id = "generate-audio-email"
  arn       = module.generate_tale_audio_lambda.lambda_function.arn
}

resource "aws_lambda_permission" "send_text_email_permission" {
  action        = "lambda:InvokeFunction"
  function_name = module.send_tale_text_lambda.lambda_function.function_name
  principal     = "events.amazonaws.com"
  source_arn    = aws_cloudwatch_event_rule.tale_text_created_event_rule.arn
}

resource "aws_lambda_permission" "generate_audio_permission" {
  action        = "lambda:InvokeFunction"
  function_name = module.generate_tale_audio_lambda.lambda_function.function_name
  principal     = "events.amazonaws.com"
  source_arn    = aws_cloudwatch_event_rule.tale_text_created_event_rule.arn
}

resource "aws_s3_bucket_notification" "notify_event_bus" {
  bucket      = aws_s3_bucket.tales.id
  eventbridge = true
}
