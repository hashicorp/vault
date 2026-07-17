/**
 * Copyright IBM Corp. 2016, 2025
 * SPDX-License-Identifier: BUSL-1.1
 */

import { service } from '@ember/service';
import Route from '@ember/routing/route';
import ListRoute from 'core/mixins/list-route';
import { paginate } from 'core/utils/paginate-list';

export default Route.extend(ListRoute, {
  pagination: service(),
  api: service(),
  version: service(),
  capabilities: service(),

  shouldReturnEmptyModel(policyType, version) {
    return policyType !== 'acl' && (version.isCommunity || !version.hasSentinel);
  },

  async model(params) {
    const policyType = this.policyType();
    if (this.shouldReturnEmptyModel(policyType, this.version)) {
      return;
    }
    try {
      const { keys } = await this.getPolicies(policyType);
      const paginated = paginate(keys, {
        page: params.page,
        filter: params.pageFilter,
      });

      const paths = paginated.map((id) =>
        this.capabilities.pathFor('policy', { policyType: this.policyType(), id })
      );
      const capabilities = paths ? await this.capabilities.fetch(paths) : {};

      const policies = paginated.map((name, i) => {
        return {
          name,
          capabilities: capabilities[paths[i]] || null,
          policyType,
        };
      });
      policies.meta = paginated.meta;
      return policies;
    } catch (err) {
      const { status } = await this.api.parseError(err);
      // acls will never be empty, but sentinel policies can be
      if (status === 404 && policyType !== 'acl') {
        return [];
      } else {
        throw err;
      }
    }
  },

  async getPolicies(policyType) {
    if (policyType === 'rgp') {
      return this.api.sys.systemListPoliciesRgp(true);
    } else if (policyType === 'egp') {
      return this.api.sys.systemListPoliciesEgp(true);
    } else {
      return this.api.sys.policiesListAclPolicies(true);
    }
  },

  setupController(controller, model) {
    const params = this.paramsFor(this.routeName);
    if (!model) {
      controller.setProperties({
        model: null,
        policyType: this.policyType(),
      });
      return;
    }
    controller.setProperties({
      model,
      filter: params.pageFilter || '',
      page: model.meta?.currentPage || 1,
      policyType: this.policyType(),
    });
  },

  resetController(controller, isExiting) {
    this._super(...arguments);
    if (isExiting) {
      controller.set('filter', '');
    }
  },

  actions: {
    willTransition(transition) {
      window.scrollTo(0, 0);
      if (!transition || transition.targetName !== this.routeName) {
        this.pagination.clearDataset();
      }
      return true;
    },
    reload() {
      this.pagination.clearDataset();
      this.refresh();
    },
  },

  policyType() {
    return this.paramsFor('vault.cluster.policies').type;
  },
});
