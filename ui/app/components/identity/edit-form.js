/**
 * Copyright IBM Corp. 2016, 2026
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
import { performSaveOperation, extractSavedId } from 'vault/utils/identity-helpers';

export default class IdentityEditFormComponent extends Component {
  @service flashMessages;
  @service api;

  @tracked policies = [];
  @tracked policyForm;
  @tracked errorBanner;

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
    const identityType = model?.identityType;
    const isAlias = model?.form?.identityFormType === 'alias';

    if (mode === 'merge') {
      return 'vault.cluster.access.identity';
    }

    if (mode === 'create') {
      return isAlias ? 'vault.cluster.access.identity.aliases' : 'vault.cluster.access.identity';
    }

    if (mode === 'edit') {
      return isAlias ? 'vault.cluster.access.identity.aliases.show' : 'vault.cluster.access.identity.show';
    }

    // Fallback route in unexpected modes.
    return identityType ? 'vault.cluster.access.identity.show' : 'vault.cluster.access.identity';
  }

  get cancelModelId() {
    const { model } = this.args;
    const isAlias = model?.form?.identityFormType === 'alias';
    if (isAlias) {
      return model?.id || model?.canonicalId;
    }
    return model?.itemId || model?.id;
  }

  getMessage(model, isDelete = false) {
    const mode = this.args.mode;
    const typeDisplay = humanize([model.identityType]);

    if (isDelete) {
      return `Successfully deleted ${typeDisplay}.`;
    }
    if (mode === 'merge') {
      return 'Successfully merged entities';
    }
    if (model.form.identityFormType === 'alias') {
      return `Successfully saved ${typeDisplay} alias.`;
    }
    const id = model.itemId || model.id;
    if (id) {
      return `Successfully saved ${typeDisplay} ${id}.`;
    }
    return `Successfully saved ${typeDisplay}.`;
  }

  save = task(
    waitFor(async () => {
      const { model, mode, onSave } = this.args;
      const { data } = model.form.toJSON();

      try {
        const response = await performSaveOperation({
          api: this.api,
          model,
          mode,
          data,
        });

        const message = this.getMessage(model);
        this.flashMessages.success(message);

        await onSave({
          saveType: 'save',
          model,
          id: extractSavedId({ mode, data, response, model }),
        });
      } catch (err) {
        const { message } = await this.api.parseError(err);
        this.errorBanner = message;
      }
    })
  );

  @action
  async deleteItem(model) {
    const message = this.getMessage(model, true);
    const flash = this.flashMessages;

    const formType = model.form.identityFormType;
    const identityType = model.identityType;

    if (formType === 'alias') {
      const methodType = identityType === 'group' ? 'groupDeleteAliasById' : 'entityDeleteAliasById';
      await this.api.identity[methodType](model.id);
    } else {
      const methodType = identityType === 'group' ? 'groupDeleteById' : 'entityDeleteById';
      await this.api.identity[methodType](model.itemId);
    }

    flash.success(message);

    return this.args.onSave({ saveType: 'delete', model });
  }

  @action
  onCreatePolicy(name) {
    this.policyForm = new PolicyForm({ name, enforcement_level: 'hard-mandatory' }, { isNew: true });
  }
}
