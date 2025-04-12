/**
 * Copyright (c) HashiCorp, Inc.
 * SPDX-License-Identifier: BUSL-1.1
 */

import Route from '@ember/routing/route';
import { service } from '@ember/service';

import type Store from '@ember-data/store';
import type SecretEngineModel from 'vault/models/secret-engine';
import type VersionService from 'vault/services/version';
import { hash } from 'rsvp';

// This route file is TODO
export default class SecretsBackendConfigurationEdit extends Route {
  @service declare readonly store: Store;
  @service declare readonly version: VersionService;

  async model() {
    const { backend } = this.paramsFor('vault.cluster.secrets.backend');
    const secretEngineRecord = this.modelFor('vault.cluster.secrets.backend') as SecretEngineModel;
    const type = secretEngineRecord.type;

    // TODO are there any mount types that we don't allow tuning?

    try {
      const response = await this.store.queryRecord('identity/oidc/config', {});
      model['identity-oidc-config'] = response;
    } catch (e) {
      // return a property called queryIssuerError and let the component handle it.
      model['identity-oidc-config'] = { queryIssuerError: true };
    }

    return hash({
      secretEngineRecord,
      backend,
      type,
    });
  }
}
