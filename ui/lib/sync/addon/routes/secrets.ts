/**
 * Copyright (c) HashiCorp, Inc.
 * SPDX-License-Identifier: BUSL-1.1
 */

import Route from '@ember/routing/route';
import { service } from '@ember/service';

import type RouterService from '@ember/routing/router-service';
import type FlagsService from 'vault/services/flags';
import type PersonaService from 'vault/services/persona';

export default class SyncSecretsRoute extends Route {
  @service declare readonly router: RouterService;
  @service declare readonly flags: FlagsService;
  @service declare readonly persona: PersonaService;

  beforeModel() {
    return this.flags.fetchActivatedFeatures();
  }

  model() {
    return {
      secretsSyncPersona: this.persona.secretsSyncPersona,
    };
  }

  afterModel() {
    if (!this.flags.secretsSyncIsActivated) {
      this.router.transitionTo('vault.cluster.sync.secrets.overview');
    }
  }
}
