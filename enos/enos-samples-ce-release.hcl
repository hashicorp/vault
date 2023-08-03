# Copyright (c) HashiCorp, Inc.
# SPDX-License-Identifier: BUSL-1.1

sample "release_ce_linux_amd64_deb" {
  attributes = global.sample_attributes

  subset "agent" {
    matrix {
      arch            = ["amd64"]
      artifact_source = ["artifactory"]
      artifact_type   = ["package"]
      consul_edition  = ["ce"]
      distro          = ["ubuntu"]
      edition         = ["ce"]
    }
  }

  subset "smoke" {
    matrix {
      arch            = ["amd64"]
      artifact_source = ["artifactory"]
      artifact_type   = ["package"]
      consul_edition  = ["ce"]
      distro          = ["ubuntu"]
      edition         = ["ce"]
    }
  }

  subset "proxy" {
    matrix {
      arch            = ["amd64"]
      artifact_source = ["artifactory"]
      artifact_type   = ["package"]
      consul_edition  = ["ce"]
      distro          = ["ubuntu"]
      edition         = ["ce"]
    }
  }

  subset "upgrade" {
    matrix {
      arch            = ["amd64"]
      artifact_source = ["artifactory"]
      artifact_type   = ["package"]
      consul_edition  = ["ce"]
      distro          = ["ubuntu"]
      edition         = ["ce"]
    }
  }
}

sample "release_ce_linux_arm64_deb" {
  attributes = global.sample_attributes

  subset "agent" {
    matrix {
      arch            = ["arm64"]
      artifact_source = ["artifactory"]
      artifact_type   = ["package"]
      consul_edition  = ["ce"]
      distro          = ["ubuntu"]
      edition         = ["ce"]
    }
  }

  subset "smoke" {
    matrix {
      arch            = ["arm64"]
      artifact_source = ["artifactory"]
      artifact_type   = ["package"]
      consul_edition  = ["ce"]
      distro          = ["ubuntu"]
      edition         = ["ce"]
    }
  }

  subset "proxy" {
    matrix {
      arch            = ["arm64"]
      artifact_source = ["artifactory"]
      artifact_type   = ["package"]
      consul_edition  = ["ce"]
      distro          = ["ubuntu"]
      edition         = ["ce"]
    }
  }

  subset "upgrade" {
    matrix {
      arch            = ["arm64"]
      artifact_source = ["artifactory"]
      artifact_type   = ["package"]
      consul_edition  = ["ce"]
      distro          = ["ubuntu"]
      edition         = ["ce"]
    }
  }
}

sample "release_ce_linux_arm64_rpm" {
  attributes = global.sample_attributes

  subset "agent" {
    matrix {
      arch            = ["arm64"]
      artifact_source = ["artifactory"]
      artifact_type   = ["package"]
      consul_edition  = ["ce"]
      distro          = ["amazon_linux", "leap", "rhel"]
      edition         = ["ce"]
    }
  }

  subset "smoke" {
    matrix {
      arch            = ["arm64"]
      artifact_source = ["artifactory"]
      artifact_type   = ["package"]
      consul_edition  = ["ce"]
      distro          = ["amazon_linux", "leap", "rhel"]
      edition         = ["ce"]
    }
  }

  subset "proxy" {
    matrix {
      arch            = ["arm64"]
      artifact_source = ["artifactory"]
      artifact_type   = ["package"]
      consul_edition  = ["ce"]
      distro          = ["amazon_linux", "leap", "rhel"]
      edition         = ["ce"]
    }
  }

  subset "upgrade" {
    matrix {
      arch            = ["arm64"]
      artifact_source = ["artifactory"]
      artifact_type   = ["package"]
      consul_edition  = ["ce"]
      distro          = ["amazon_linux", "leap", "rhel"]
      edition         = ["ce"]
    }
  }
}

sample "release_ce_linux_amd64_rpm" {
  attributes = global.sample_attributes

  subset "agent" {
    matrix {
      arch            = ["amd64"]
      artifact_source = ["artifactory"]
      artifact_type   = ["package"]
      consul_edition  = ["ce"]
      distro          = ["amazon_linux", "leap", "rhel", "sles"]
      edition         = ["ce"]
    }
  }

  subset "smoke" {
    matrix {
      arch            = ["amd64"]
      artifact_source = ["artifactory"]
      artifact_type   = ["package"]
      consul_edition  = ["ce"]
      distro          = ["amazon_linux", "leap", "rhel", "sles"]
      edition         = ["ce"]
    }
  }

  subset "proxy" {
    matrix {
      arch            = ["amd64"]
      artifact_source = ["artifactory"]
      artifact_type   = ["package"]
      consul_edition  = ["ce"]
      distro          = ["amazon_linux", "leap", "rhel", "sles"]
      edition         = ["ce"]
    }
  }

  subset "upgrade" {
    matrix {
      arch            = ["amd64"]
      artifact_source = ["artifactory"]
      artifact_type   = ["package"]
      consul_edition  = ["ce"]
      distro          = ["amazon_linux", "leap", "rhel", "sles"]
      edition         = ["ce"]
    }
  }
}

sample "release_ce_linux_amd64_zip" {
  attributes = global.sample_attributes

  subset "agent" {
    matrix {
      arch            = ["amd64"]
      artifact_type   = ["bundle"]
      artifact_source = ["artifactory"]
      consul_edition  = ["ce"]
      edition         = ["ce"]
    }
  }

  subset "smoke" {
    matrix {
      arch            = ["amd64"]
      artifact_type   = ["bundle"]
      artifact_source = ["artifactory"]
      consul_edition  = ["ce"]
      edition         = ["ce"]
    }
  }

  subset "proxy" {
    matrix {
      arch            = ["amd64"]
      artifact_type   = ["bundle"]
      artifact_source = ["artifactory"]
      consul_edition  = ["ce"]
      edition         = ["ce"]
    }
  }

  subset "upgrade" {
    matrix {
      arch            = ["amd64"]
      artifact_type   = ["bundle"]
      artifact_source = ["artifactory"]
      consul_edition  = ["ce"]
      edition         = ["ce"]
    }
  }
}

sample "release_ce_linux_arm64_zip" {
  attributes = global.sample_attributes

  subset "agent" {
    matrix {
      arch            = ["arm64"]
      artifact_source = ["artifactory"]
      artifact_type   = ["bundle"]
      consul_edition  = ["ce"]
      edition         = ["ce"]
    }
  }

  subset "smoke" {
    matrix {
      arch            = ["arm64"]
      artifact_source = ["artifactory"]
      artifact_type   = ["bundle"]
      edition         = ["ce"]
    }
  }

  subset "proxy" {
    matrix {
      arch            = ["arm64"]
      artifact_source = ["artifactory"]
      artifact_type   = ["bundle"]
      consul_edition  = ["ce"]
      edition         = ["ce"]
    }
  }

  subset "upgrade" {
    matrix {
      arch            = ["arm64"]
      artifact_source = ["artifactory"]
      artifact_type   = ["bundle"]
      consul_edition  = ["ce"]
      edition         = ["ce"]
    }
  }
}
