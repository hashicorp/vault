/**
 * Copyright (c) HashiCorp, Inc.
 * SPDX-License-Identifier: BUSL-1.1
 */

import Route from '@ember/routing/route';
import { service } from '@ember/service';
import SecretsEngineForm from 'vault/forms/secrets/engine';

import type { ModelFrom } from 'vault/vault/route';
import type Store from '@ember-data/store';

export type MountSecretBackendModel = ModelFrom<VaultClusterSettingsMountSecretBackendRoute>;

export default class VaultClusterSettingsMountSecretBackendRoute extends Route {
  @service declare readonly store: Store;

  model() {
    const defaults = {
      config: { listingVisibility: false },
      kvConfig: {
        maxVersions: 0,
        casRequired: false,
        deleteVersionAfter: undefined,
      },
      options: { version: 2 },
    };
    return new SecretsEngineForm(defaults, { isNew: true });
  }
}
