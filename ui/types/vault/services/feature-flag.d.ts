/**
 * Copyright (c) HashiCorp, Inc.
 * SPDX-License-Identifier: BUSL-1.1
 */

import Service from '@ember/service';

export default class FeatureFlagService extends Service {
  featureFlags: string[] | null;
  setFeatureFlags: () => void;
  managedNamespaceRoot: string | null;
}
