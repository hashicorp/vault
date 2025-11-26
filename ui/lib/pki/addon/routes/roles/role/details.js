/**
 * Copyright IBM Corp. 2016, 2025
 * SPDX-License-Identifier: BUSL-1.1
 */

import Route from '@ember/routing/route';
import { service } from '@ember/service';

export default class RolesRoleDetailsRoute extends Route {
  @service api;
  @service secretMountPath;
  @service capabilities;

  async fetchCapabilities(id) {
    const { pathFor } = this.capabilities;
    const backend = this.secretMountPath.currentPath;

    const pathMap = {
      role: pathFor('pkiRole', { backend, id }),
      issue: pathFor('pkiIssue', { backend, id }),
      sign: pathFor('pkiSign', { backend, id }),
    };
    const perms = await this.capabilities.fetch(Object.values(pathMap));

    return {
      canEdit: perms[pathMap.role].canUpdate,
      canDelete: perms[pathMap.role].canDelete,
      canGenerateCert: perms[pathMap.issue].canUpdate,
      canSign: perms[pathMap.sign].canUpdate,
    };
  }

  async model() {
    const { role: name } = this.paramsFor('roles/role');
    return {
      role: await this.api.secrets
        .pkiReadRole(name, this.secretMountPath.currentPath)
        .then((role) => ({ name, ...role })),
      capabilities: await this.fetchCapabilities(name),
    };
  }

  setupController(controller, resolvedModel) {
    super.setupController(controller, resolvedModel);
    controller.breadcrumbs = [
      { label: 'Secrets', route: 'secrets', linkExternal: true },
      { label: this.secretMountPath.currentPath, route: 'overview', model: this.secretMountPath.currentPath },
      { label: 'Roles', route: 'roles.index', model: this.secretMountPath.currentPath },
      { label: resolvedModel.role.name },
    ];
  }
}
