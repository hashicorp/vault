/**
 * Copyright (c) HashiCorp, Inc.
 * SPDX-License-Identifier: BUSL-1.1
 */

import { service } from '@ember/service';
import Component from '@glimmer/component';
import { task } from 'ember-concurrency';
import { waitFor } from '@ember/test-waiters';
import { tracked } from '@glimmer/tracking';

import type AuthMethodForm from 'vault/forms/auth/method';
import type RouterService from '@ember/routing/router-service';
import type FlashMessageService from 'ember-cli-flash/services/flash-messages';
import type ApiService from 'vault/services/api';
import type { HTMLElementEvent } from 'vault/forms';
import {
  AwsConfigureClientRequest,
  AwsConfigureIdentityAccessListTidyOperationRequest,
  AwsConfigureRoleTagDenyListTidyOperationRequest,
  AzureConfigureAuthRequest,
  GithubConfigureRequest,
  GoogleCloudConfigureAuthRequest,
  JwtConfigureRequest,
  KubernetesConfigureAuthRequest,
  LdapConfigureAuthRequest,
  OktaConfigureRequest,
  RadiusConfigureRequest,
} from '@hashicorp/vault-client-typescript';
import AuthMethodResource from 'vault/resources/auth/method';

/**
 * @module AuthConfigForm/Config
 * The `AuthConfigForm/Config` is the form for auth methods that need additional configuration.
 * AuthConfigForm::Options handle the backend's mount configuration.
 *
 * @example
 * <AuthConfigForm::Config @form={{this.form}} />
 *
 * @property form=null {AuthMethodForm} - The corresponding auth method that is being configured.
 *
 */

type ConfigPayload =
  | AwsConfigureClientRequest
  | AwsConfigureIdentityAccessListTidyOperationRequest
  | AwsConfigureRoleTagDenyListTidyOperationRequest
  | AzureConfigureAuthRequest
  | GithubConfigureRequest
  | GoogleCloudConfigureAuthRequest
  | JwtConfigureRequest
  | KubernetesConfigureAuthRequest
  | LdapConfigureAuthRequest
  | OktaConfigureRequest
  | RadiusConfigureRequest
  | Record<string, unknown>; // Add other payload types as needed

interface Args {
  form: AuthMethodForm;
  section: 'configuration' | 'client' | 'identity-accesslist' | 'roletag-denylist';
  method: AuthMethodResource;
}

export default class AuthConfigBase extends Component<Args> {
  @service declare readonly api: ApiService;
  @service declare readonly flashMessages: FlashMessageService;
  @service declare readonly router: RouterService;

  @tracked errorMessage = '';

  configMethod(path: string, payload: ConfigPayload) {
    const { section, method } = this.args;

    switch (method.methodType) {
      case 'aws':
        switch (section) {
          case 'client':
            return this.api.auth.awsConfigureClient(path, payload as AwsConfigureClientRequest);
          case 'identity-accesslist':
            return this.api.auth.awsConfigureIdentityAccessListTidyOperation(
              path,
              payload as AwsConfigureIdentityAccessListTidyOperationRequest
            );
          case 'roletag-denylist':
            return this.api.auth.awsConfigureRoleTagDenyListTidyOperation(
              path,
              payload as AwsConfigureRoleTagDenyListTidyOperationRequest
            );
          default:
            throw new Error(`Unsupported AWS section: ${section}`);
        }
      case 'azure':
        return this.api.auth.azureConfigureAuth(path, payload as AzureConfigureAuthRequest);
      case 'github':
        return this.api.auth.githubConfigure(path, payload as GithubConfigureRequest);
      case 'gcp':
        return this.api.auth.googleCloudConfigureAuth(path, payload as GoogleCloudConfigureAuthRequest);
      case 'jwt':
      case 'oidc':
        return this.api.auth.jwtConfigure(path, payload as JwtConfigureRequest);
      case 'kubernetes':
        return this.api.auth.kubernetesConfigureAuth(path, payload as KubernetesConfigureAuthRequest);
      case 'ldap':
        return this.api.auth.ldapConfigureAuth(path, payload as LdapConfigureAuthRequest);
      case 'okta':
        return this.api.auth.oktaConfigure(path, payload as OktaConfigureRequest);
      case 'radius':
        return this.api.auth.radiusConfigure(path, payload as RadiusConfigureRequest);
      default:
        throw new Error(`Configuration of the ${method.methodType} method is not supported by the Vault UI.`);
    }
  }

  @task
  @waitFor
  *saveModel(evt: HTMLElementEvent<HTMLFormElement>) {
    evt.preventDefault();
    this.errorMessage = '';
    try {
      const { form, method } = this.args;
      const { data } = form.toJSON();
      yield this.configMethod(method.path, data as ConfigPayload);
      this.router.transitionTo('vault.cluster.access.methods').followRedirects();
      this.flashMessages.success('The configuration was saved successfully.');
    } catch (err) {
      const { message } = yield this.api.parseError(err);
      this.errorMessage = message;
    }
  }
}
