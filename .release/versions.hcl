# Copyright (c) HashiCorp, Inc.
# SPDX-License-Identifier: BUSL-1.1

# This manifest file describes active branches and is consumed by various
# pipeline tooling. It ought to be considered the source of truth regarding
# which branches are currently active, whether the Ent version is considered
# LTS, and whether or not a branch's CE counterpart is active. Main is always
# assumed to be active but any version not present will be assumed inactive.

schema = 1
active_versions {
  version "1.20.x" {
    ce_active = true
  }

  version "1.19.x" {
    ce_active = true
    lts       = true
  }

  version "1.18.x" {
    ce_active = false
  }

  version "1.17.x" {
    ce_active = false
  }

  version "1.16.x" {
    ce_active = false
    lts       = true
  }
}
