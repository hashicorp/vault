/**
 * Copyright IBM Corp. 2016, 2025
 * SPDX-License-Identifier: BUSL-1.1
 */

import Controller from '@ember/controller';
import { action } from '@ember/object';
import { service } from '@ember/service';

export default class PolicyEditController extends Controller {
  @service router;
  @service api;
  @service flashMessages;

  @action
  async deletePolicy() {
    const { policyType, name } = this.model;
    try {
      if (policyType === 'acl') {
        await this.api.sys.policiesDeleteAclPolicy(name);
      } else if (policyType === 'egp') {
        await this.api.sys.systemDeletePoliciesEgpName(name);
      } else {
        await this.api.sys.systemDeletePoliciesRgpName(name);
      }

      this.flashMessages.success(`${policyType.toUpperCase()} policy "${name}" was successfully deleted.`);
      this.router.transitionTo('vault.cluster.policies', policyType);
    } catch (error) {
      const { status, message } = await this.api.parseError(error);
      if (status === 404) {
        return [];
      }
      const flashMessage = `There was an error deleting the ${policyType.toUpperCase()} policy "${name}": ${message}.`;
      this.flashMessages.danger(flashMessage);
    }
  }
}
