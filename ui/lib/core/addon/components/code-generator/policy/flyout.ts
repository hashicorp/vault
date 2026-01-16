/**
 * Copyright IBM Corp. 2016, 2025
 * SPDX-License-Identifier: BUSL-1.1
 */

import { action } from '@ember/object';
import Component from '@glimmer/component';
import { tracked } from '@glimmer/tracking';
import { policySnippetArgs, PolicyStanza } from 'core/utils/code-generators/policy';
import { validate } from 'vault/utils/forms/validate';
import { service } from '@ember/service';
import { task } from 'ember-concurrency';

import type { HTMLElementEvent } from 'vault/forms';
import type { PolicyData } from './builder';
import type { ValidationMap, Validations } from 'vault/vault/app-types';
import type ApiService from 'vault/services/api';
import type FlashMessageService from 'ember-cli-flash/services/flash-messages';

export default class CodeGeneratorPolicyFlyout extends Component {
  @service declare readonly api: ApiService;
  @service declare readonly flashMessages: FlashMessageService;

  @tracked errorMessage = '';
  @tracked policyContent = '';
  @tracked policyName = '';
  @tracked showFlyout = false;
  @tracked stanzas: PolicyStanza[] = [new PolicyStanza()];
  @tracked validationErrors: ValidationMap | null = null;

  validations: Validations = {
    name: [{ type: 'presence', message: 'Name is required.' }],
  };

  validationError = (param: string) => {
    const { isValid, errors } = this.validationErrors?.[param] ?? {};
    return !isValid && errors ? errors.join(' ') : '';
  };

  @action
  handleNameInput(event: HTMLElementEvent<HTMLInputElement>) {
    const { value } = event.target;
    this.policyName = value.toLowerCase();
  }

  @action
  handlePolicyChange({ policy, stanzas }: PolicyData) {
    this.policyContent = policy;
    this.stanzas = stanzas;
  }

  get snippetArgs() {
    const policyName = this.policyName || '<policy name>';
    return policySnippetArgs(policyName, this.stanzas);
  }

  @task
  *onSave(event: HTMLElementEvent<HTMLFormElement>) {
    event.preventDefault();
    this.resetErrors();

    const { isValid, state, invalidFormMessage } = validate({ name: this.policyName }, this.validations);
    if (!isValid) {
      this.validationErrors = state;
      this.errorMessage = invalidFormMessage;
      return;
    }

    try {
      yield this.api.sys.policiesWriteAclPolicy(this.policyName, { policy: this.policyContent });
      this.flashMessages.success(`ACL policy "${this.policyName}" saved successfully.`, {
        link: {
          text: 'View policy',
          route: 'vault.cluster.policy.show',
          models: ['acl', this.policyName],
        },
      });
      this.showFlyout = false;
      this.resetFlyoutState();
    } catch (e) {
      const { message } = yield this.api.parseError(e);
      this.errorMessage = message;
    }
  }

  @action
  onClose() {
    this.showFlyout = false;
    this.resetErrors();
  }

  resetErrors() {
    this.validationErrors = null;
    this.errorMessage = '';
  }

  resetFlyoutState() {
    this.policyName = '';
    this.policyContent = '';
    this.stanzas = [new PolicyStanza()];
  }
}
