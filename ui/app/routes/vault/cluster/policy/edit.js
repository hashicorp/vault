/**
 * Copyright IBM Corp. 2016, 2026
 * SPDX-License-Identifier: BUSL-1.1
 */

import Route from '@ember/routing/route';
import { service } from '@ember/service';
import PolicyForm from 'vault/forms/policy';

export default class PolicyEditRouter extends Route {
  @service router;
  @service api;
  @service capabilities;

  beforeModel() {
    const params = this.paramsFor(this.routeName);
    const policyType = this.policyType();
    if (policyType === 'acl' && params.policy_name === 'root') {
      return this.router.transitionTo('vault.cluster.policies.acl');
    }
  }

  async model(params) {
    // Reuse the show route's model only when it matches the policy being edited.
    // show.currentModel persists on the route singleton after navigation, so
    // checking the name prevents a stale model from a previously-visited policy
    // being served to the wrong edit form (VAULT-50826).
    const showModel = this.modelFor('vault.cluster.policy.show');
    if (showModel && showModel.name === params.policy_name) {
      const form = new PolicyForm(showModel, { isNew: false });
      form.policyType = showModel.policyType;
      form.capabilities = showModel.capabilities;
      return form;
    }

    const type = this.policyType();
    const policy = await this.fetchPolicy(params.policy_name, type);
    const form = new PolicyForm(policy, { isNew: false });
    form.policyType = type;
    form.capabilities = await this.capabilities.for('policy', {
      policyType: this.policyType(),
      id: params.policy_name,
    });
    return form;
  }

  async fetchPolicy(name, type) {
    let res;
    if (type === 'acl') {
      res = await this.api.sys.policiesReadAclPolicy(name);
    } else if (type === 'egp') {
      res = (await this.api.sys.systemReadPoliciesEgpName(name)).data;
    } else {
      res = (await this.api.sys.systemReadPoliciesRgpName(name)).data;
    }

    return { name, ...res };
  }

  policyType() {
    return this.paramsFor('vault.cluster.policy').type;
  }
}
