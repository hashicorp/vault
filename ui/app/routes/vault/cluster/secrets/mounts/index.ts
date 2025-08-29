/**
 * Copyright (c) HashiCorp, Inc.
 * SPDX-License-Identifier: BUSL-1.1
 */

import Route from '@ember/routing/route';
import { service } from '@ember/service';
import SecretsEngineForm from 'vault/forms/secrets/engine';
import Router from 'vault/router';

import type { ModelFrom } from 'vault/vault/route';

export type MountSecretBackendModel = ModelFrom<VaultClusterSecretsMountsIndexRouter>;

export default class VaultClusterSecretsMountsIndexRouter extends Route {
  @service declare router: Router;

  model() {
    const defaults = {
      config: { listing_visibility: false },
      kv_config: {
        max_versions: 0,
        cas_required: false,
        delete_version_after: undefined,
      },
      options: { version: 2 },
    };
    return new SecretsEngineForm(defaults, { isNew: true });
  }
}
