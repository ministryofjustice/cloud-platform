provider "aws" {
  region = "eu-west-1"
}

data "aws_iam_policy_document" "asb-cfn" {
  statement {
    actions = [
      "cloudformation:*",
      "iam:*",
      "kms:*",
      "ssm:*",
      "ec2:*",
      "lambda:*",
      "athena:*",
      "dynamodb:*",
      "elasticache:*",
      "elasticmapreduce:*",
      "rds:*",
      "redshift:*",
      "route53:*",
      "s3:*",
      "sns:*",
      "sqs:*",
    ]

    resources = [
      "*",
    ]
  }
}

resource "aws_iam_policy" "asb-cfn" {
  name        = "${terraform.workspace}_asb-cfn"
  description = "${terraform.workspace} AWS Service Broken Cloudformation Policy"
  path        = "/system/asb/${terraform.workspace}/"
  policy      = "${data.aws_iam_policy_document.asb-cfn.json}"
}

data "aws_iam_policy_document" "asb-cfn-trust" {
  statement {
    actions = ["sts:AssumeRole"]

    principals {
      type        = "Service"
      identifiers = ["cloudformation.amazonaws.com"]
    }
  }
}

resource "aws_iam_role" "asb-cfn" {
  name               = "${terraform.workspace}_asb-cfn"
  description        = "${terraform.workspace}: AWS Service Broker Cloudformation Role"
  path               = "/system/asb/${terraform.workspace}/"
  assume_role_policy = "${data.aws_iam_policy_document.asb-cfn-trust.json}"
}

resource "aws_iam_policy_attachment" "asb-cfn" {
  name       = "deploy"
  roles      = ["${aws_iam_role.asb-cfn.name}"]
  policy_arn = "${aws_iam_policy.asb-cfn.arn}"
}

data "aws_iam_policy_document" "asb" {
  statement {
    actions = [
      "cloudformation:CancelUpdateStack",
      "cloudformation:ContinueUpdateRollback",
      "cloudformation:CreateStack",
      "cloudformation:CreateUploadBucket",
      "cloudformation:DeleteStack",
      "cloudformation:DescribeAccountLimits",
      "cloudformation:DescribeStackEvents",
      "cloudformation:DescribeStackResource",
      "cloudformation:DescribeStackResources",
      "cloudformation:DescribeStacks",
      "cloudformation:GetStackPolicy",
      "cloudformation:ListStackResources",
      "cloudformation:ListStacks",
      "cloudformation:SetStackPolicy",
      "cloudformation:UpdateStack",
      "iam:AddUserToGroup",
      "iam:AttachUserPolicy",
      "iam:CreateAccessKey",
      "iam:CreatePolicy",
      "iam:CreatePolicyVersion",
      "iam:CreateUser",
      "iam:DeleteAccessKey",
      "iam:DeletePolicy",
      "iam:DeletePolicyVersion",
      "iam:DeleteRole",
      "iam:DeleteUser",
      "iam:DeleteUserPolicy",
      "iam:DetachUserPolicy",
      "iam:GetPolicy",
      "iam:GetPolicyVersion",
      "iam:GetUser",
      "iam:GetUserPolicy",
      "iam:ListAccessKeys",
      "iam:ListGroups",
      "iam:ListGroupsForUser",
      "iam:ListInstanceProfiles",
      "iam:ListPolicies",
      "iam:ListPolicyVersions",
      "iam:ListRoles",
      "iam:ListUserPolicies",
      "iam:ListUsers",
      "iam:PutUserPolicy",
      "iam:RemoveUserFromGroup",
      "iam:UpdateUser",
      "ec2:DescribeVpcs",
      "ec2:DescribeSubnets",
      "ec2:DescribeAvailabilityZones",
    ]

    resources = [
      "*",
    ]
  }

  statement {
    actions = [
      "iam:PassRole",
    ]

    resources = [
      "${aws_iam_role.asb-cfn.arn}",
    ]
  }

  statement {
    actions = [
      "iam:PassRole",
    ]

    resources = [
      "${aws_iam_role.asb-cfn.arn}",
    ]
  }

  // do we need this?
  statement {
    actions = [
      "ssm:GetParameters",
    ]

    resources = [
      "arn:aws:ssm:*:*:parameter/asb-access-key-id-*",
      "arn:aws:ssm:*:*:parameter/asb-secret-access-key-*",
    ]
  }
}

resource "aws_iam_policy" "asb" {
  name        = "${terraform.workspace}_asb"
  description = "${terraform.workspace} AWS Service Broken Policy"
  path        = "/system/asb/${terraform.workspace}/"
  policy      = "${data.aws_iam_policy_document.asb.json}"
}

resource "aws_iam_user" "asb" {
  name = "${terraform.workspace}_asb"
  path = "/system/asb/${terraform.workspace}/"
}

resource "aws_iam_policy_attachment" "asb" {
  name       = "asb"
  users      = ["${aws_iam_user.asb.name}"]
  policy_arn = "${aws_iam_policy.asb.arn}"
}

resource "aws_iam_access_key" "asb" {
  user = "${aws_iam_user.asb.name}"
}

output "AWS_ACCESS_KEY_ID" {
  value = "${aws_iam_access_key.asb.id}"
}

output "AWS_SECRET_ACCESS_KEY" {
  value = "${aws_iam_access_key.asb.secret}"
}

output "CFN_ROLE_ARN" {
  value = "${aws_iam_role.asb-cfn.arn}"
}
