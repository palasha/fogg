# Auto-generated by fogg. Do not edit
# Make improvements in fogg, so that everyone can benefit.


provider "aws" {
  version             = "~> 0.15.0"
  region              = "us-west-env1"
  profile             = "czi-env"
  allowed_account_ids = [5]
}

# Aliased Providers (for doing things in every region).


provider "aws" {
  alias               = "us-east-env2"
  version             = "~> 0.15.0"
  region              = "us-east-env2"
  profile             = "czi-env"
  allowed_account_ids = [5]
}














terraform {
  required_version = "~>0.15.0"

  backend "s3" {
    bucket         = "env-bucket"
    dynamodb_table = "env-table"

    key = "terraform/env-project/envs/stage/components/cloud-env.tfstate"


    encrypt = true
    region  = "us-west-env1"
    profile = "czi-env"
  }
}

variable "env" {
  type    = string
  default = "stage"
}

variable "project" {
  type    = string
  default = "env-project"
}


variable "region" {
  type    = string
  default = "us-west-env1"
}


variable "component" {
  type    = string
  default = "cloud-env"
}


variable "aws_profile" {
  type    = string
  default = "czi-env"
}



variable "owner" {
  type    = string
  default = "env@example.com"
}

variable "tags" {
  type = map(string)
  default = {
    project   = "env-project"
    env       = "stage"
    service   = "cloud-env"
    owner     = "env@example.com"
    managedBy = "terraform"
  }
}


variable "foo" {
  type    = string
  default = "env"
}


data "terraform_remote_state" "global" {
  backend = "s3"

  config = {
    bucket         = "the-bucket"
    dynamodb_table = "the-table"
    key            = "terraform/test-project/global.tfstate"
    region         = "us-west-2"
    profile        = "czi"
  }
}



data "terraform_remote_state" "helm" {
  backend = "s3"

  config = {
    bucket         = "env-bucket"
    dynamodb_table = "env-table"
    key            = "terraform/env-project/envs/stage/components/helm.tfstate"
    region         = "us-west-env1"
    profile        = "czi-env"
  }
}


# remote state for accounts

data "terraform_remote_state" "bar" {
  backend = "s3"

  config = {
    bucket         = "bar-bucket"
    dynamodb_table = "bar-table"
    key            = "terraform/bar-project/accounts/bar.tfstate"
    region         = "us-west-bar1"
    profile        = "czi-bar"
  }
}

data "terraform_remote_state" "foo" {
  backend = "s3"

  config = {
    bucket         = "foo-bucket"
    dynamodb_table = "foo-table"
    key            = "terraform/foo-project/accounts/foo.tfstate"
    region         = "us-west-foo1"
    profile        = "czi-foo"
  }
}


# map of aws_accounts
variable "aws_accounts" {
  type = map
  default = {


    bar = 3



    foo = 2


  }
}

provider random {
  version = "~> 2.2"
}

provider template {
  version = "~> 2.1"
}

provider archive {
  version = "~> 1.3"
}

provider null {
  version = "~> 2.1"
}

provider local {
  version = "~> 1.4"
}

provider tls {
  version = "~> 2.1"
}
