/**
 * Copyright IBM Corp. 2016, 2025
 * SPDX-License-Identifier: BUSL-1.1
 */

import Route from '@ember/routing/route';
import { service } from '@ember/service';
import { ModelFrom } from 'vault/route';

import type { KubernetesApplicationModel } from '../application';
import type Store from '@ember-data/store';
import type SecretMountPath from 'vault/services/secret-mount-path';
import type Controller from '@ember/controller';
import type { Breadcrumb } from 'vault/app-types';
import type Transition from '@ember/routing/transition';
import type AdapterError from 'vault/@ember-data/adapter/error';

interface RouteController extends Controller {
  breadcrumbs: Array<Breadcrumb>;
  model: KubernetesRolesModel;
}

export type KubernetesRolesModel = ModelFrom<KubernetesRolesRoute>;

export default class KubernetesRolesRoute extends Route {
  @service declare readonly store: Store;
  @service declare readonly secretMountPath: SecretMountPath;

  async model(_params: unknown, transition: Transition) {
    const { promptConfig, secretsEngine } = this.modelFor('application') as KubernetesApplicationModel;
    const model = { promptConfig, secretsEngine, roles: [] };
    try {
      // filter roles based on pageFilter value
      const { pageFilter } = (transition.to?.queryParams || {}) as { pageFilter?: string };
      const models = await this.store.query('kubernetes/role', { backend: this.secretMountPath.currentPath });
      const roles = pageFilter
        ? models.filter((model) => model.name.toLowerCase().includes(pageFilter.toLowerCase()))
        : models;
      return {
        ...model,
        roles: roles as unknown[],
      };
    } catch (error) {
      if ((error as AdapterError).httpStatus !== 404) {
        throw error;
      }
    }

    return model;
  }

  setupController(controller: RouteController, resolvedModel: KubernetesRolesModel) {
    super.setupController(controller, resolvedModel);
    const { currentPath } = this.secretMountPath;
    controller.breadcrumbs = [
      { label: 'Secrets', route: 'secrets', linkExternal: true },
      { label: currentPath, route: 'overview', model: currentPath },
      { label: 'Roles' },
    ];
  }
}
