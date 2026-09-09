/**
 * Copyright IBM Corp. 2016, 2025
 * SPDX-License-Identifier: BUSL-1.1
 */

import Route from '@ember/routing/route';
import UnloadModelRoute from 'vault/mixins/unload-model-route';
import { service } from '@ember/service';

/**
 * @type Class
 */
export default Route.extend(UnloadModelRoute, {
  router: service(),
  api: service(),
  capabilities: service(),

  beforeModel() {
    const params = this.paramsFor(this.routeName);
    const policyType = this.policyType();
    if (policyType === 'acl' && params.policy_name === 'root') {
      return this.router.transitionTo('vault.cluster.policies', 'acl');
    }
  },

  async model(params) {
    const type = this.policyType();

    let res;
    if (type === 'acl') {
      res = await this.api.sys.policiesReadAclPolicy(params.policy_name);
    } else if (type === 'egp') {
      res = (await this.api.sys.systemReadPoliciesEgpName(params.policy_name)).data;
    } else {
      res = (await this.api.sys.systemReadPoliciesRgpName(params.policy_name)).data;
    }

    const policy = { name: params.policy_name, ...res };

    return {
      ...policy,
      policyType: type,
      format: this.format(res.policy),
      capabilities: await this.capabilities.for('policy', {
        policyType: this.policyType(),
        id: params.policy_name,
      }),
    };
  },

  format(policy) {
    let isJSON;
    try {
      const parsed = JSON.parse(policy);
      if (parsed) {
        isJSON = true;
      }
    } catch (e) {
      // can't parse JSON
      isJSON = false;
    }
    return isJSON ? 'json' : 'hcl';
  },

  policyType() {
    return this.paramsFor('vault.cluster.policy').type;
  },
});
