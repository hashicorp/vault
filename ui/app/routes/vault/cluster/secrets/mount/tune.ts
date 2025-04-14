/**
 * Copyright (c) HashiCorp, Inc.
 * SPDX-License-Identifier: BUSL-1.1
 */

import Route from '@ember/routing/route';
import { service } from '@ember/service';

import type Store from '@ember-data/store';
import ApiService from 'vault/services/api';
import type SecretEngineModel from 'vault/models/secret-engine';
import type VersionService from 'vault/services/version';
import { hash } from 'rsvp';

import type { Breadcrumb } from 'vault/vault/app-types';
import type Controller from '@ember/controller';
import type Transition from '@ember/routing/transition';

interface RouteModel {
  secretEngineRecord: SecretEngineModel;
  backend: string;
}

interface RouteController extends Controller {
  breadcrumbs: Array<Breadcrumb>;
  model: RouteModel;
}

// This route file is TODO
export default class SecretsBackendConfigurationEdit extends Route {
  @service declare readonly api: ApiService;
  @service declare readonly store: Store;
  @service declare readonly version: VersionService;

  async model(params: { mount_name: string }) {
    const mount_name = params.mount_name;
    const secretEngineModel = this.store.createRecord('secret-engine');
    secretEngineModel.set('path', mount_name);
    return secretEngineModel;
  }

  setupController(controller: RouteController, resolvedModel: RouteModel, transition: Transition) {
    super.setupController(controller, resolvedModel, transition);

    const crumbs = [
      { label: 'Secrets', route: 'vault.cluster.secrets' },
      { label: resolvedModel.backend, route: 'vault.cluster.secrets.backend' },
      { label: 'Configuration', route: 'vault.cluster.secrets.backend.configuration' },
      { label: 'Tune mount' },
    ];

    controller.set('breadcrumbs', crumbs);
  }
}
