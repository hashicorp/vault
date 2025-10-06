/**
 * Copyright (c) HashiCorp, Inc.
 * SPDX-License-Identifier: BUSL-1.1
 */

import Route from '@ember/routing/route';
import SecretsEngineForm from 'vault/forms/secrets/engine';

export default class VaultClusterSecretsMountsCreateRoute extends Route {
  model(params: { mount_type: string }) {
    const { mount_type } = params;

    const defaults = {
      path: mount_type, // Default path to match the engine type
      config: { listing_visibility: false },
      kv_config: {
        max_versions: 0,
        cas_required: false,
        delete_version_after: undefined,
      },
      options: { version: 2 },
    };

    const form = new SecretsEngineForm(defaults, { isNew: true });
    // Explicitly set the type on the form after creation
    form.type = mount_type;
    // Apply type-specific defaults (e.g., PKI max lease TTL)
    form.applyTypeSpecificDefaults();

    return form;
  }
}
