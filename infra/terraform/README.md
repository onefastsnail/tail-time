# Terraform

## Setup

An account per environment

### 1. Setup CI user

- Setup CI user group
- Setup CI user in CI user group
- Create policy and attach to CI user group

```json
{
  "Version": "2012-10-17",
  "Statement": [
    {
      "Sid": "VisualEditor0",
      "Effect": "Allow",
      "Action": [
        "ses:*",
        "events:*",
        "s3:*",
        "logs:*",
        "iam:*",
        "cloudfront:*",
        "cloudwatch:*",
        "kms:*",
        "lambda:*"
      ],
      "Resource": "*"
    }
  ]
}
```
