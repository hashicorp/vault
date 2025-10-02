/**
 * Copyright (c) HashiCorp, Inc.
 * SPDX-License-Identifier: BUSL-1.1
 */

import Component from '@glimmer/component';
import { service } from '@ember/service';
import { task } from 'ember-concurrency';
import { waitFor } from '@ember/test-waiters';
import { tracked } from '@glimmer/tracking';
import { supportedTypes } from 'vault/utils/auth-form-helpers';

import type AuthMethodForm from 'vault/forms/auth/method';
import type ApiService from 'vault/services/api';
import type FlashMessageService from 'vault/services/flash-messages';
import type RouterService from '@ember/routing/router-service';
import type { HTMLElementEvent } from 'vault/forms';
import type NamespaceService from 'vault/services/namespace';
import type VersionService from 'vault/services/version';
import type { MountsAuthTuneConfigurationParametersRequest } from '@hashicorp/vault-client-typescript';

/**
 * @module AuthConfigForm/Options
 * The `AuthConfigForm/Options` is options portion of the auth config form.
 *
 * @example
 * <AuthConfigForm::Options @form={{this.form}} />
 *
 * @property form=null {AuthMethodForm} - The corresponding auth method that is being configured.
 *
 */

type Args = {
  form: AuthMethodForm;
};

export default class AuthConfigOptions extends Component<Args> {
  @service declare readonly api: ApiService;
  @service declare readonly flashMessages: FlashMessageService;
  @service declare readonly router: RouterService;
  @service declare readonly namespace: NamespaceService;
  @service declare readonly version: VersionService;

  @tracked errorMessage: string | null = null;

  get directLoginLink() {
    const ns = this.namespace.path;
    const nsQueryParam = ns ? `namespace=${encodeURIComponent(ns)}&` : '';
    const { normalizedType, data } = this.args.form;
    const isSupported = supportedTypes(this.version.isEnterprise).includes(normalizedType);
    return isSupported
      ? `${window.origin}/ui/vault/auth?${nsQueryParam}with=${encodeURIComponent(data.path)}`
      : '';
  }

  get supportsUserLockoutConfig() {
    return ['approle', 'ldap', 'userpass'].includes(this.args.form.normalizedType);
  }

  onSubmit = task(
    waitFor(async (evt: HTMLElementEvent<HTMLFormElement>) => {
      evt.preventDefault();
      this.errorMessage = null;
      try {
        const { form } = this.args;
        const {
          data: { description, config, user_lockout_config },
        } = form.toJSON();

        const payload = {
          description,
          ...config,
        } as MountsAuthTuneConfigurationParametersRequest;

        if (Object.keys(user_lockout_config).length) {
          payload.user_lockout_config = user_lockout_config;
        }

        // 'token_type' cannot be set for the 'token' auth mount
        if (form.normalizedType === 'token' && payload?.token_type) {
          delete payload.token_type;
        }

        await this.api.sys.mountsAuthTuneConfigurationParameters(form.data.path, payload);
        this.router.transitionTo('vault.cluster.access.methods').followRedirects();
        this.flashMessages.success('The configuration was saved successfully.');
      } catch (err) {
        const { message } = await this.api.parseError(err);
        this.errorMessage = message;
      }
    })
  );
}
