terraform {
  extra_arguments "private_vars" {
    commands = [
      "init",
      "apply",
      "refresh",
      "import",
      "plan",
      "taint",
      "untaint",
      "destroy"
    ]
  }
}

generate "provider" {
  path = "provider.tf"
  if_exists = "overwrite_terragrunt"
  contents = <<EOF
terraform {
  required_providers {
    device42 = {
      version = "~> 0.2.0"
      source  = "github.com/chopnico/device42"
    }
  }
}

provider "device42" {
  username   = yamldecode(file("../../private.yml"))["username"]
  password   = yamldecode(file("../../private.yml"))["password"]
  host       = yamldecode(file("../../private.yml"))["host"]
  proxy	     = "http://127.0.0.1:8080"
  ignore_ssl = true
}
EOF
}
