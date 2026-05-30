/**
 * Copyright IBM Corp. 2016, 2025
 * SPDX-License-Identifier: BUSL-1.1
 */

import { action } from '@ember/object';
import Component from '@glimmer/component';
import { tracked } from '@glimmer/tracking';
import { formatStanzas, policySnippetArgs, PolicyStanza } from 'core/utils/code-generators/policy';
import { validate } from 'vault/utils/forms/validate';
import { service } from '@ember/service';
import { task } from 'ember-concurrency';
import routerLookup from 'core/utils/router-lookup';

import type { HTMLElementEvent } from 'vault/forms';
import type { PolicyData } from './builder';
import type { ValidationMap, Validations } from 'vault/vault/app-types';
import type ApiService from 'vault/services/api';
import type CapabilitiesService from 'vault/services/capabilities';
import type FlashMessageService from 'ember-cli-flash/services/flash-messages';
import type VersionService from 'vault/services/version';
import type RouterService from '@ember/routing/router-service';

interface Args {
  onClose?: CallableFunction;
  policyPaths: string[];
}

export default class CodeGeneratorPolicyFlyout extends Component<Args> {
  @service declare readonly api: ApiService;
  @service declare readonly capabilities: CapabilitiesService;
  @service declare readonly flashMessages: FlashMessageService;
  @service declare readonly version: VersionService;

  defaultStanzas = [new PolicyStanza()];

  validations: Validations = {
    name: [{ type: 'presence', message: 'Name is required.' }],
    stanzas: [
      {
        validator: ({ stanzas }) =>
          stanzas.length > 0 && stanzas.every((stanza: PolicyStanza) => stanza.isValid),
        message: 'Invalid policy content.',
      },
    ],
  };

  @tracked errorDetails: string[] = [];
  @tracked errorMessage = '';
  @tracked policyName = '';
  @tracked showFlyout = false;
  @tracked stanzas: PolicyStanza[] = this.defaultStanzas;
  @tracked validationErrors: ValidationMap | null = null;

  get router(): RouterService {
    return routerLookup(this);
  }

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
  handlePolicyChange({ stanzas }: PolicyData) {
    this.stanzas = stanzas;
  }

  get policyContent() {
    return formatStanzas(this.stanzas);
  }

  get snippetArgs() {
    const policyName = this.policyName || '<policy name>';
    const policy = formatStanzas(this.stanzas);
    return policySnippetArgs(policyName, policy);
  }

  @task
  *onSave(event: HTMLElementEvent<HTMLFormElement>) {
    event.preventDefault();
    this.resetErrors();

    const { isValid, state } = validate({ name: this.policyName, stanzas: this.stanzas }, this.validations);

    if (!isValid) {
      this.validationErrors = state;
      this.errorDetails = Object.values(state).flatMap((s) => s.errors);
      // Render general error message instead of exact count from validate() because
      // stanzas (which are validated as a single input) can have up to 2 errors each.
      const msg = this.errorDetails.length > 1 ? 'are errors' : 'is an error';
      this.errorMessage = `There ${msg} with this form.`;
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
      this.closeFlyout();
      this.resetFlyoutState();
    } catch (e) {
      const { message } = yield this.api.parseError(e);
      this.errorMessage = message;
    }
  }

  @action
  openFlyout() {
    this.showFlyout = true;

    const defaultState = formatStanzas(this.defaultStanzas);
    const currentState = formatStanzas(this.stanzas);
    const noChanges = currentState === defaultState;
    // Only preset stanzas if no changes have been made to the flyout
    if (noChanges) {
      const presetStanzas = this.computePolicyPaths();
      this.stanzas = presetStanzas ? presetStanzas : this.stanzas;
    }
  }

  @action
  closeFlyout() {
    this.showFlyout = false;
    this.resetErrors();
    this.args.onClose?.();
  }

  resetErrors() {
    this.validationErrors = null;
    this.errorMessage = '';
    this.errorDetails = [];
  }

  resetFlyoutState() {
    this.policyName = '';
    this.stanzas = [new PolicyStanza()];
  }

  computePolicyPaths() {
    // Explicit policy paths take precedence
    if (this.args.policyPaths) {
      return this.presetStanzas(this.args.policyPaths);
    }

    const currentRouteName = this.router.currentRouteName;
    if (!currentRouteName) {
      return null;
    }

    const routePolicyPaths = this.capabilities.lookupRoutePaths(currentRouteName);
    return routePolicyPaths ? this.presetStanzas(Array.from(routePolicyPaths)) : null;
  }

  presetStanzas(paths: string[]) {
    const presetStanzas = paths.map((path) => new PolicyStanza({ path }));
    return presetStanzas.length ? presetStanzas : null;
  }
}
