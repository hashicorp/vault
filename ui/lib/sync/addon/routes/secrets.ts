/**
 * Copyright (c) HashiCorp, Inc.
 * SPDX-License-Identifier: BUSL-1.1
 */

import Route from '@ember/routing/route';
import { service } from '@ember/service';

import type RouterService from '@ember/routing/router-service';
import type StoreService from 'vault/services/store';
import type VersionService from 'vault/services/version';
import type PersonaService from 'vault/services/persona';

// interface ActivationFlagsResponse {
//   data: {
//     activated: Array<string>;
//     unactivated: Array<string>;
//   };
// }

export default class SyncSecretsRoute extends Route {
  @service declare readonly router: RouterService;
  @service declare readonly store: StoreService;
  @service declare readonly version: VersionService;
  @service declare readonly persona: PersonaService;

  beforeModel() {
    return this.version.fetchActivatedFeatures();
  }

  model() {
    return {
      secretsSyncPersona: this.persona.secretsSyncPersona,
    };
  }

  // ARG TODO return, what was this doing?
  // afterModel(model: { activatedFeatures: Array<string> }) {
  //   if (!model.activatedFeatures) {
  //     this.router.transitionTo('vault.cluster.sync.secrets.overview');
  //   }
  // }
}
