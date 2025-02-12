/**
 * Copyright (c) HashiCorp, Inc.
 * SPDX-License-Identifier: BUSL-1.1
 */

import { service } from '@ember/service';
import Route from '@ember/routing/route';
import { pluralize } from 'ember-inflector';
import { capitalize } from '@ember/string';
import ListRoute from 'vault/mixins/list-route';

export default Route.extend(ListRoute, {
  pagination: service(),
  pathHelp: service('path-help'),
  api: service(),

  getMethodAndModelInfo() {
    const { item_type: itemType } = this.paramsFor('vault.cluster.access.method.item');
    const { path: authMethodPath } = this.paramsFor('vault.cluster.access.method');
    const methodModel = this.modelFor('vault.cluster.access.method');
    const { apiPath, type } = methodModel;
    return { apiPath, type, authMethodPath, itemType, methodModel };
  },

  async model() {
    const { type, authMethodPath, itemType } = this.getMethodAndModelInfo();
    const { page, pageFilter } = this.paramsFor(this.routeName);
    const payload = {
      [`${type}MountPath`]: authMethodPath,
      list: true,
    };
    // examples -> userpassListUser, kubernetesListAuthRoles, ldapListGroups
    const listItem = type === 'kubernetes' && itemType === 'role' ? 'authRole' : itemType;
    const authListMethod = `${type}List${capitalize(pluralize(listItem))}`;

    try {
      const { keys } = await this.api.auth[authListMethod](payload);
      // it would likely be better to update the template/component to use the keys directly
      // for now we are trying to make as few changes as possible
      const mappedKeys = keys.map((key) => ({ id: key }));
      return this.pagination.paginate(mappedKeys, {
        page,
        pageSize: 3,
        filter: pageFilter,
        filterKey: 'id',
      });
    } catch (error) {
      const err = (await error.response?.json()) || error;
      if (err.httpStatus === 404) {
        return [];
      } else {
        throw err;
      }
    }
  },

  actions: {
    willTransition(transition) {
      window.scrollTo(0, 0);
      if (transition.targetName !== this.routeName) {
        this.pagination.clearDataset();
      }
      return true;
    },
    reload() {
      this.pagination.clearDataset();
      this.refresh();
    },
  },

  setupController(controller) {
    this._super(...arguments);
    const { apiPath, authMethodPath, itemType, methodModel } = this.getMethodAndModelInfo();
    controller.set('itemType', itemType);
    controller.set('methodModel', methodModel);
    this.pathHelp.getPaths(apiPath, authMethodPath, itemType).then((paths) => {
      controller.set(
        'paths',
        paths.paths.filter((path) => path.navigation && path.itemType.includes(itemType))
      );
    });
  },
});
