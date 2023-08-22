terraform {
  required_providers {
    device42 = {
      version = "~> 0.2.0"
      source  = "github.com/chopnico/device42"
    }
  }
}

provider "device42" {
  username = yamldecode(file("private.yml"))["username"]
  password = yamldecode(file("private.yml"))["password"]
  host     = yamldecode(file("private.yml"))["host"]

  # you may need to ignore ssl errors if your web certificate is
  # not publicly trusted or its certificate chain is not install
  # on your system.
  ignore_ssl = true

  # this was originally created to test http/s transactions
  # using mitmproxy, but it could be useful if you are required
  # to use a proxy.
  proxy = "http://127.0.0.1:8080"
}
