/**
 * Copyright IBM Corp. 2016, 2026
 * SPDX-License-Identifier: BUSL-1.1
 */

import Route from '@ember/routing/route';
import { service } from '@ember/service';

import type Controller from '@ember/controller';
import type SecretMountPath from 'vault/services/secret-mount-path';
import type { Breadcrumb } from 'vault/app-types';

interface RouteController extends Controller {
  title: string;
  breadcrumbs: Array<Breadcrumb>;
}

export default class PkiExternalErrorRoute extends Route {
  @service declare readonly secretMountPath: SecretMountPath;

  setupController(controller: RouteController, resolvedModel: unknown) {
    super.setupController(controller, resolvedModel);
    const { currentPath } = this.secretMountPath;
    controller.title = currentPath;
    controller.breadcrumbs = [
      { label: 'Vault', route: 'vault', icon: 'vault', linkExternal: true },
      { label: 'Secrets engines', route: 'secrets', linkExternal: true },
      { label: 'Public PKI' },
    ];
  }
}
