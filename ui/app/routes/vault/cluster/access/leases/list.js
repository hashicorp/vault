/**
 * Copyright IBM Corp. 2016, 2026
 * SPDX-License-Identifier: BUSL-1.1
 */

import { hash } from 'rsvp';
import Route from '@ember/routing/route';
import { action } from '@ember/object';
import { service } from '@ember/service';
import { paginate } from 'core/utils/paginate-list';
import { keyIsFolder, keyWithoutParentKey } from 'core/utils/key-utils';
import { SystemApiLeasesLookUpListEnum } from '@hashicorp/vault-client-typescript';

export default class LeasesListRoute extends Route {
  @service api;
  @service capabilities;

  queryParams = {
    page: {
      refreshModel: true,
    },
    pageFilter: {
      refreshModel: true,
    },
  };

  templateName = 'vault/cluster/access/leases/list';

  async model(params) {
    const prefix = params.prefix || '';

    if (!this.modelFor('vault.cluster.access.leases').canList) {
      return undefined;
    }

    const leasesResult = await this.api.sys
      .leasesLookUp(prefix, SystemApiLeasesLookUpListEnum.TRUE)
      .then((resp) => {
        const rawKeys = resp.keys ?? [];
        return rawKeys.map((key) => {
          const id = prefix ? prefix + key : key;
          return {
            id,
            isFolder: keyIsFolder(id),
            keyWithoutParent: keyWithoutParentKey(id),
          };
        });
      })
      .catch((err) => {
        if (err.status === 404 && prefix === '') {
          return [];
        }
        throw err;
      });

    const leases = paginate(leasesResult, {
      page: params.page ? Number(params.page) : 1,
      filter: params.pageFilter || undefined,
    });

    const capabilities = await hash({
      revokePrefix: this.capabilities.fetchPathCapabilities(`sys/leases/revoke-prefix/${prefix}`),
      forceRevokePrefix: this.capabilities.fetchPathCapabilities(`sys/leases/revoke-force/${prefix}`),
    });

    return { leases, capabilities };
  }

  setupController(controller, model) {
    const params = this.paramsFor(this.routeName);
    const prefix = params.prefix ? params.prefix : '';
    controller.set('hasModel', true);
    controller.setProperties({
      model: model?.leases,
      capabilities: model?.capabilities,
      baseKey: { id: prefix },
    });
    if (model?.leases) {
      const pageFilter = params.pageFilter;
      let filter;
      if (prefix) {
        filter = prefix + (pageFilter || '');
      } else if (pageFilter) {
        filter = pageFilter;
      }
      controller.setProperties({
        filter: filter || '',
        page: model.leases?.meta?.currentPage,
      });
    }
  }

  resetController(controller, isExiting) {
    super.resetController(...arguments);
    if (isExiting) {
      controller.set('filter', '');
    }
  }

  @action
  error(error, transition) {
    const { prefix } = this.paramsFor(this.routeName);
    error.keyId = prefix;
    // eslint-disable-next-line ember/no-controller-access-in-routes
    const hasModel = this.controllerFor(this.routeName).hasModel;
    if (hasModel && error.status === 404) {
      transition.abort();
    } else {
      return true;
    }
  }

  @action
  willTransition() {
    window.scrollTo(0, 0);
    return true;
  }
}
