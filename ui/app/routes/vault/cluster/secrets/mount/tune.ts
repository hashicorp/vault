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

  async model() {
    const { backend } = this.paramsFor('vault.cluster.secrets.backend');
    const secretEngineRecord = this.modelFor('vault.cluster.secrets.backend') as SecretEngineModel;
    const type = secretEngineRecord.type;
    let tuneMountResponse;

    // TODO are there any mount types that we don't allow tuning?

    try {
      // const response = await this.api.sys.mountsReadTuningInformation(
      //   { path: secretEngineRecord.id },
      //   this.api.buildHeaders({ token: '' })
      // );
      tuneMountResponse = {
        default_lease_ttl: 3600,
        max_lease_ttl: 7200,
        force_no_cache: false,
        plugin_version: 'some-string', // plugin version does not get returned on read?
      };
    } catch (e) {
      // todo handle error
    }

    return hash({
      secretEngineRecord,
      backend,
      type,
      tuneMount: tuneMountResponse,
    });
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
