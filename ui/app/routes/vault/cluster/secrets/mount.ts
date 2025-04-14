/**
 * Copyright (c) HashiCorp, Inc.
 * SPDX-License-Identifier: BUSL-1.1
 */

import Route from '@ember/routing/route';
import { service } from '@ember/service';

import type { ModelFrom } from 'vault/vault/route';
import type Store from '@ember-data/store';
import type RouterService from '@ember/routing/router-service';

import VaultClusterSettingsMountSecretBackendRoute from '../settings/mount-secret-backend';
import type { Breadcrumb } from 'vault/vault/app-types';
import type Controller from '@ember/controller';
import type Transition from '@ember/routing/transition';
import type SecretEngineModel from 'vault/models/secret-engine';
import { hash } from 'rsvp';

export type MountSecretBackendModel = ModelFrom<VaultClusterSettingsMountSecretBackendRoute>; // ARG TODO change this after update

interface RouteModel {
  secretEngineRecord: SecretEngineModel;
  mountConfigRecord: SecretEngineModel; // ARG todo break this apart.
  backend: string;
}

interface RouteController extends Controller {
  breadcrumbs: Array<Breadcrumb>;
  model: RouteModel;
}

export default class VaultClusterSecretsMount extends Route {
  @service declare readonly store: Store;
  @service declare readonly router: RouterService;

  model() {
    const secretEngineModel = this.store.createRecord('secret-engine');
    const mountConfigModel = this.store.createRecord('mount-config');
    return hash({
      secretEngineModel,
      mountConfigModel,
    });
  }

  setupController(controller: RouteController, resolvedModel: RouteModel, transition: Transition) {
    super.setupController(controller, resolvedModel, transition);

    const crumbs = [{ label: 'Secrets', route: 'vault.cluster.secrets' }, { label: 'Mounts' }];
    // likely have to do some kind of magic with will Transition
    controller.set('breadcrumbs', crumbs);
    controller.set('router', this.router);
  }
}
