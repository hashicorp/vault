# Copyright (c) HashiCorp, Inc.
# SPDX-License-Identifier: MPL-2.0

sample "release_ce_linux_amd64_deb" {
  attributes = global.sample_attributes

  subset "agent" {
    matrix {
      arch            = ["amd64"]
      artifact_source = ["artifactory"]
      artifact_type   = ["package"]
      distro          = ["ubuntu"]
      edition         = ["ce"]
    }
  }

  subset "smoke" {
    matrix {
      arch            = ["amd64"]
      artifact_source = ["artifactory"]
      artifact_type   = ["package"]
      distro          = ["ubuntu"]
      edition         = ["ce"]
    }
  }

  subset "proxy" {
    matrix {
      arch            = ["amd64"]
      artifact_source = ["artifactory"]
      artifact_type   = ["package"]
      distro          = ["ubuntu"]
      edition         = ["ce"]
    }
  }

  subset "upgrade" {
    matrix {
      arch            = ["amd64"]
      artifact_source = ["artifactory"]
      artifact_type   = ["package"]
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
      distro          = ["ubuntu"]
      edition         = ["ce"]
    }
  }

  subset "smoke" {
    matrix {
      arch            = ["arm64"]
      artifact_source = ["artifactory"]
      artifact_type   = ["package"]
      distro          = ["ubuntu"]
      edition         = ["ce"]
    }
  }

  subset "proxy" {
    matrix {
      arch            = ["arm64"]
      artifact_source = ["artifactory"]
      artifact_type   = ["package"]
      distro          = ["ubuntu"]
      edition         = ["ce"]
    }
  }

  subset "upgrade" {
    matrix {
      arch            = ["arm64"]
      artifact_source = ["artifactory"]
      artifact_type   = ["package"]
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
      distro          = ["rhel"]
      edition         = ["ce"]
    }
  }

  subset "smoke" {
    matrix {
      arch            = ["arm64"]
      artifact_source = ["artifactory"]
      artifact_type   = ["package"]
      distro          = ["rhel"]
      edition         = ["ce"]
    }
  }

  subset "proxy" {
    matrix {
      arch            = ["arm64"]
      artifact_source = ["artifactory"]
      artifact_type   = ["package"]
      distro          = ["rhel"]
      edition         = ["ce"]
    }
  }

  subset "upgrade" {
    matrix {
      arch            = ["arm64"]
      artifact_source = ["artifactory"]
      artifact_type   = ["package"]
      distro          = ["rhel"]
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
      distro          = ["rhel"]
      edition         = ["ce"]
    }
  }

  subset "smoke" {
    matrix {
      arch            = ["amd64"]
      artifact_source = ["artifactory"]
      artifact_type   = ["package"]
      distro          = ["rhel"]
      edition         = ["ce"]
    }
  }

  subset "proxy" {
    matrix {
      arch            = ["amd64"]
      artifact_source = ["artifactory"]
      artifact_type   = ["package"]
      distro          = ["rhel"]
      edition         = ["ce"]
    }
  }

  subset "upgrade" {
    matrix {
      arch            = ["amd64"]
      artifact_source = ["artifactory"]
      artifact_type   = ["package"]
      distro          = ["rhel"]
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
      edition         = ["ce"]
    }
  }

  subset "smoke" {
    matrix {
      arch            = ["amd64"]
      artifact_type   = ["bundle"]
      artifact_source = ["artifactory"]
      edition         = ["ce"]
    }
  }

  subset "proxy" {
    matrix {
      arch            = ["amd64"]
      artifact_type   = ["bundle"]
      artifact_source = ["artifactory"]
      edition         = ["ce"]
    }
  }

  subset "upgrade" {
    matrix {
      arch            = ["amd64"]
      artifact_type   = ["bundle"]
      artifact_source = ["artifactory"]
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
      edition         = ["ce"]
    }
  }

  subset "upgrade" {
    matrix {
      arch            = ["arm64"]
      artifact_source = ["artifactory"]
      artifact_type   = ["bundle"]
      edition         = ["ce"]
    }
  }
}
