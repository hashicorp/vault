schema = "1"

project "vault" {
  team = "vault"
  slack {
    notification_channel = "#feed-releng" #TODO update slack channel
  }
  github {
    organization = "hashicorp"
    repository = "vault"
    release_branches = ["release/1.8.x"]
  }
}

event "merge" {
  // "entrypoint" to use if build is not run automatically
  // i.e. send "merge" complete signal to orchestrator to trigger build
}

event "build" {
  depends = ["merge"]
  action "build" {
    organization = "hashicorp"
    repository = "vault"
    workflow = "build"
  }
}

event "upload-dev" {
  depends = ["build"]
  action "upload-dev" {
    organization = "hashicorp"
    repository = "crt-workflows-common"
    workflow = "upload-dev"
    depends = ["build"]
  }

  notification {
    on = "fail"
  }
}

event "quality-tests" {
  depends = ["upload-dev"]
  action "quality-tests" {
    organization = "hashicorp"
    repository = "crt-workflows-common"
    workflow = "quality-tests"
  }

  notification {
    on = "fail"
  }
}

event "security-scan" {
  depends = ["quality-tests"]
  action "security-scan" {
    organization = "hashicorp"
    repository = "crt-workflows-common"
    workflow = "security-scan"
  }

  notification {
    on = "fail"
  }
}

event "notarize-darwin-amd64" {
  depends = ["security-scan"]
  action "notarize-darwin-amd64" {
    organization = "hashicorp"
    repository = "crt-workflows-common"
    workflow = "notarize-darwin-amd64"
  }

  notification {
    on = "fail"
  }
}

event "notarize-darwin-arm64" {
  depends = ["notarize-darwin-amd64"]
  action "notarize-darwin-arm64" {
    organization = "hashicorp"
    repository = "crt-workflows-common"
    workflow = "notarize-darwin-arm64"
  }

  notification {
    on = "fail"
  }
}

event "notarize-windows-386" {
  depends = ["notarize-darwin-arm64"]
  action "notarize-windows-386" {
    organization = "hashicorp"
    repository = "crt-workflows-common"
    workflow = "notarize-windows-386"
  }

  notification {
    on = "fail"
  }
}

event "notarize-windows-amd64" {
  depends = ["notarize-windows-386"]
  action "notarize-windows-amd64" {
    organization = "hashicorp"
    repository = "crt-workflows-common"
    workflow = "notarize-windows-amd64"
  }

  notification {
    on = "fail"
  }
}

event "sign" {
  depends = ["notarize-windows-amd64"]
  action "sign" {
    organization = "hashicorp"
    repository = "crt-workflows-common"
    workflow = "sign"
  }

  notification {
    on = "fail"
  }
}

event "sign-linux-rpms" {
  depends = ["sign"]
  action "sign-linux-rpms" {
    organization = "hashicorp"
    repository = "crt-workflows-common"
    workflow = "sign-linux-rpms"
  }

  notification {
    on = "fail"
  }
}

event "verify" {
  depends = ["sign-linux-rpms"]
  action "verify" {
    organization = "hashicorp"
    repository = "crt-workflows-common"
    workflow = "verify"
  }

  notification {
    on = "fail"
  }
}

event "promote-staging" {

  action "promote-staging" {
    organization = "hashicorp"
    repository = "crt-workflows-common"
    workflow = "promote-staging"
  }

  notification {
    on = "fail"
  }

  notification {
    on = "success"
  }
}

event "promote-production" {

  action "promote-production" {
    organization = "hashicorp"
    repository = "crt-workflows-common"
    workflow = "promote-production"
  }

  notification {
    on = "fail"
  }

  notification {
    on = "success"
  }
}

event "post-publish" {
  depends = ["promote-production"]

  action "post-publish" {
    organization = "hashicorp"
    repository = "crt-workflows-common"
    workflow = "post-publish"
  }

  notification {
    on = "fail"
  }

  notification {
    on = "success"
  }
}
