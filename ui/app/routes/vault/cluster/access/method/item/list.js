/**
 * Copyright (c) HashiCorp, Inc.
 * SPDX-License-Identifier: BUSL-1.1
 */

import { service } from '@ember/service';
import Route from '@ember/routing/route';
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

  fetchListItems(type, itemType, authMethodPath) {
    const { auth } = this.api;
    if (type === 'userpass') {
      return auth.userpassListUsers({ userpassMountPath: authMethodPath, list: 'true' });
    }
    if (type === 'kubernetes') {
      return auth.kubernetesListAuthRoles({ kubernetesMountPath: authMethodPath, list: 'true' });
    }
    if (type === 'ldap') {
      const payload = {
        ldapMountPath: authMethodPath,
        list: 'true' /* as const */,
      };
      return itemType === 'group' ? auth.ldapListGroups(payload) : auth.ldapListUsers(payload);
    }
    if (type === 'okta') {
      const payload = {
        oktaMountPath: authMethodPath,
        list: 'true' /* as const */,
      };
      return itemType === 'group' ? auth.oktaListGroups(payload) : auth.oktaListUsers(payload);
    }
    if (type === 'radius') {
      return auth.radiusListUsers({ radiusMountPath: authMethodPath, list: 'true' });
    }
  },

  async model() {
    const { type, authMethodPath, itemType } = this.getMethodAndModelInfo();
    const { page, pageFilter } = this.paramsFor(this.routeName);
    // examples -> userpassListUser, kubernetesListAuthRoles, ldapListGroups
    const listItem = type === 'kubernetes' && itemType === 'role' ? 'authRole' : itemType;

    try {
      const { keys } = await this.fetchListItems(type, itemType, authMethodPath);
      // it would likely be better to update the template/component to use the keys directly
      // for now we are trying to make as few changes as possible
      // add some additional information necessary to delete method in generated-item-list component
      const mappedKeys = keys.map((key) => ({
        id: key,
        type,
        listItem,
        authMethodPath,
      }));
      return this.pagination.paginate(mappedKeys, {
        page,
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
    const { itemType, methodModel } = this.getMethodAndModelInfo();
    controller.set('itemType', itemType);
    controller.set('methodModel', methodModel);
  },
});
