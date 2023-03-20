# Terraform

## Setup

An account per environment.

## Prerequisites

- Install [Terraform](https://developer.hashicorp.com/terraform/downloads)

### 1. Setup AWS CI user and permissions

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
        "cloudwatch:*",
        "kms:*",
        "lambda:*"
      ],
      "Resource": "*"
    }
  ]
}
```

### 2. Setup the environment

Copy and configure `.env.sample` to `.env` and use it.

### 3. Terraform

- `terraform workspace select dev` to activate a workspace per environment.   
- `terraform init` to install all needed Terraform dependencies.
- `terraform plan -out=infra.tfplan` to plan what will be provisioned.
- `terraform apply infra.tfplan` to apply the plan created which compiles and deploys the apps and infra.
