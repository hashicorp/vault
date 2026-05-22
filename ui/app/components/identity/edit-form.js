/**
 * Copyright IBM Corp. 2016, 2025
 * SPDX-License-Identifier: BUSL-1.1
 */

import Component from '@glimmer/component';
import { tracked } from '@glimmer/tracking';
import { service } from '@ember/service';
import { action } from '@ember/object';
import { task } from 'ember-concurrency';
import { humanize } from 'vault/helpers/humanize';
import { waitFor } from '@ember/test-waiters';
import PolicyForm from 'vault/forms/policy';
import {
  SystemApiPoliciesListAclPoliciesListEnum,
  SystemApiSystemListPoliciesRgpListEnum,
} from '@hashicorp/vault-client-typescript';

export default class IdentityEditFormComponent extends Component {
  @service flashMessages;
  @service store;
  @service api;

  @tracked policies = [];
  @tracked policyForm;

  constructor() {
    super(...arguments);
    // fetch policies to populate dropdown in form
    this.fetchPolicies();
  }

  async fetchPolicies() {
    const [aclResult, rgpResult] = await Promise.allSettled([
      this.api.sys.policiesListAclPolicies(SystemApiPoliciesListAclPoliciesListEnum.TRUE),
      this.api.sys.systemListPoliciesRgp(SystemApiSystemListPoliciesRgpListEnum.TRUE),
    ]);
    const aclPolicies = aclResult.status === 'fulfilled' ? aclResult.value.keys : [];
    const rgpPolicies = rgpResult.status === 'fulfilled' ? rgpResult.value.keys : [];
    this.policies = [...aclPolicies, ...rgpPolicies]
      .filter((name) => name !== 'root')
      .map((policy) => ({ id: policy }));
  }

  get cancelLink() {
    const { model, mode } = this.args;
    const routes = {
      'create-entity': 'vault.cluster.access.identity',
      'edit-entity': 'vault.cluster.access.identity.show',
      'merge-entity-merge': 'vault.cluster.access.identity',
      'create-entity-alias': 'vault.cluster.access.identity.aliases',
      'edit-entity-alias': 'vault.cluster.access.identity.aliases.show',
      'create-group': 'vault.cluster.access.identity',
      'edit-group': 'vault.cluster.access.identity.show',
      'create-group-alias': 'vault.cluster.access.identity.aliases',
      'edit-group-alias': 'vault.cluster.access.identity.aliases.show',
    };
    const key = model ? `${mode}-${model.identityType}` : 'merge-entity-alias';
    return routes[key];
  }

  getMessage(model, isDelete = false) {
    const mode = this.mode;
    const typeDisplay = humanize([model.identityType]);
    const action = isDelete ? 'deleted' : 'saved';
    if (mode === 'merge') {
      return 'Successfully merged entities';
    }
    if (model.id) {
      return `Successfully ${action} ${typeDisplay} ${model.id}.`;
    }
    return `Successfully ${action} ${typeDisplay}.`;
  }

  save = task(
    waitFor(async () => {
      const { model } = this.args;
      const message = this.getMessage(model);

      try {
        await model.save();
      } catch (err) {
        // err will display via model state
        return;
      }
      this.flashMessages.success(message);
      await this.args.onSave({ saveType: 'save', model });
    })
  );

  willDestroy() {
    // components are torn down after store is disconnected and will cause an error if attempt to unload record
    const noTeardown = this.store && !this.store.isDestroying;
    const model = this.model;
    if (noTeardown && model && model.isDirty && !model.isDestroyed && !model.isDestroying) {
      model.rollbackAttributes();
    }
    super.willDestroy(...arguments);
  }

  @action
  deleteItem(model) {
    const message = this.getMessage(model, true);
    const flash = this.flashMessages;
    model.destroyRecord().then(() => {
      flash.success(message);
      return this.args.onSave({ saveType: 'delete', model });
    });
  }
  @action
  onCreatePolicy(name) {
    this.policyForm = new PolicyForm({ name, enforcement_level: 'hard-mandatory' }, { isNew: true });
  }
}
