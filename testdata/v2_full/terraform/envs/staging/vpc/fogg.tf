# Auto-generated by fogg. Do not edit
# Make improvements in fogg, so that everyone can benefit.


provider "aws" {
  version             = "~> 0.12.0"
  region              = "us-west-2"
  profile             = "profile"
  allowed_account_ids = [456]
}

# Aliased Providers (for doing things in every region).













terraform {
  required_version = "~>0.100.0"

  backend "s3" {
    bucket = "buck"


    key = "terraform/proj/envs/staging/components/vpc.tfstate"


    encrypt = true
    region  = "us-west-2"
    profile = "profile"
  }
}

variable "env" {
  type    = "string"
  default = "staging"
}

variable "project" {
  type    = "string"
  default = "proj"
}


variable "region" {
  type    = "string"
  default = "us-west-2"
}


variable "component" {
  type    = "string"
  default = "vpc"
}


variable "aws_profile" {
  type    = "string"
  default = "profile"
}



variable "owner" {
  type    = "string"
  default = "foo@example.com"
}

variable "tags" {
  type = "map"
  default = {
    project   = "proj"
    env       = "staging"
    service   = "vpc"
    owner     = "foo@example.com"
    managedBy = "terraform"
  }
}


variable "foo" {
  type    = "string"
  default = "bar3"
}


data "terraform_remote_state" "global" {
  backend = "s3"

  config = {
    bucket = "buck"

    key     = "terraform/proj/global.tfstate"
    region  = "us-west-2"
    profile = "profile"
  }
}



data "terraform_remote_state" "comp1" {
  backend = "s3"

  config = {
    bucket = "buck"

    key     = "terraform/proj/envs/staging/components/comp1.tfstate"
    region  = "us-west-2"
    profile = "profile"
  }
}

data "terraform_remote_state" "comp2" {
  backend = "s3"

  config = {
    bucket = "buck"

    key     = "terraform/proj/envs/staging/components/comp2.tfstate"
    region  = "us-west-2"
    profile = "profile"
  }
}


# remote state for accounts

data "terraform_remote_state" "bar" {
  backend = "s3"

  config = {
    bucket = "buck"

    key     = "terraform/proj/accounts/bar.tfstate"
    region  = "us-west-2"
    profile = "profile"
  }
}

data "terraform_remote_state" "foo" {
  backend = "s3"

  config = {
    bucket = "buck"

    key     = "terraform/proj/accounts/foo.tfstate"
    region  = "us-west-2"
    profile = "profile"
  }
}


# map of aws_accounts
variable "aws_accounts" {
  type = "map"
  default = {


    bar = 456



    foo = 123


  }
}
