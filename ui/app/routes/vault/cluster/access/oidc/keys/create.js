/**
 * Copyright IBM Corp. 2016, 2025
 * SPDX-License-Identifier: BUSL-1.1
 */

import Route from '@ember/routing/route';
import { service } from '@ember/service';
import OidcKeyForm from 'vault/forms/oidc/key';

export default class OidcKeysCreateRoute extends Route {
  @service api;

  model() {
    const defaultValues = {
      algorithm: 'RS256',
      rotation_period: '24h',
      verification_ttl: '24h',
    };
    return new OidcKeyForm(defaultValues, { isNew: true });
  }
}
