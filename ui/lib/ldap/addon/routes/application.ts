/**
 * Copyright IBM Corp. 2016, 2025
 * SPDX-License-Identifier: BUSL-1.1
 */

import Route from '@ember/routing/route';
import { service } from '@ember/service';

import type { ModelFrom } from 'vault/vault/route';
import type ApiService from 'vault/services/api';
import type Transition from '@ember/routing/transition';
import type SecretsEngineResource from 'vault/resources/secrets/engine';
import type { LdapConfigureRequest } from '@hashicorp/vault-client-typescript';

export type LdapApplicationModel = ModelFrom<LdapApplicationRoute>;

export default class LdapApplicationRoute extends Route {
  @service declare readonly api: ApiService;

  async model(params: Record<string, unknown>, transition: Transition) {
    const secretsEngine = super.model(params, transition) as SecretsEngineResource;
    let config: LdapConfigureRequest | undefined;
    let promptConfig = false;
    let configError: unknown;
    // check if engine is configured
    // child routes will handle prompting for configuration if needed
    try {
      const { data } = await this.api.secrets.ldapReadConfiguration(secretsEngine.id);
      config = data as LdapConfigureRequest;
    } catch (error) {
      const { response, status } = await this.api.parseError(error);
      // not considering 404 an error since it triggers the cta
      if (status === 404) {
        promptConfig = true;
      } else {
        // ignore if the user does not have permission or other failures so as to not block the other operations
        // this error is thrown in the configuration route so we can display the error in the view
        configError = response;
      }
    }
    return {
      secretsEngine,
      config,
      configError,
      promptConfig,
    };
  }
}
