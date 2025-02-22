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
    if (type === 'userpass') {
      return this.api.get('/auth/{userpass_mount_path}/users/', {
        params: {
          path: { userpass_mount_path: authMethodPath },
          query: { list: 'true' },
        },
      });
    }
    if (type === 'kubernetes') {
      return this.api.get('/auth/{kubernetes_mount_path}/role/', {
        params: {
          path: { kubernetes_mount_path: authMethodPath },
          query: { list: 'true' },
        },
      });
    }
    if (type === 'ldap') {
      if (itemType === 'group') {
        return this.api.get('/auth/{ldap_mount_path}/groups/', {
          params: {
            path: { ldap_mount_path: authMethodPath },
            query: { list: 'true' },
          },
        });
      }
      if (itemType === 'user') {
        return this.api.get('/auth/{ldap_mount_path}/users/', {
          params: {
            path: { ldap_mount_path: authMethodPath },
            query: { list: 'true' },
          },
        });
      }
    }
    if (type === 'okta') {
      if (itemType === 'group') {
        return this.api.get('/auth/{okta_mount_path}/groups/', {
          params: {
            path: { okta_mount_path: authMethodPath },
            query: { list: 'true' },
          },
        });
      }
      if (itemType === 'user') {
        return this.api.get('/auth/{okta_mount_path}/users/', {
          params: {
            path: { okta_mount_path: authMethodPath },
            query: { list: 'true' },
          },
        });
      }
    }
    if (type === 'radius') {
      return this.api.get('/auth/{radius_mount_path}/users/', {
        params: {
          path: { radius_mount_path: authMethodPath },
          query: { list: 'true' },
        },
      });
    }
  },

  async model() {
    const { type, authMethodPath, itemType } = this.getMethodAndModelInfo();
    const { page, pageFilter } = this.paramsFor(this.routeName);

    const { data, error } = await this.fetchListItems(type, itemType, authMethodPath);

    if (!error) {
      // it would likely be better to update the template/component to use the keys directly
      // for now we are trying to make as few changes as possible
      const mappedKeys = data.keys.map((key) => ({ id: key }));
      return this.pagination.paginate(mappedKeys, {
        page,
        pageSize: 3,
        filter: pageFilter,
        filterKey: 'id',
      });
    }

    if (error.httpStatus === 404) {
      return [];
    }

    throw error;
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
