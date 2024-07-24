/**
 * Copyright (c) HashiCorp, Inc.
 * SPDX-License-Identifier: BUSL-1.1
 */

import { service } from '@ember/service';
import Component from '@glimmer/component';
import errorMessage from 'vault/utils/error-message';
import { tracked } from '@glimmer/tracking';

import type Store from '@ember-data/store';
import type SecretEngineModel from 'vault/models/secret-engine';
import type AdapterError from '@ember-data/adapter';

/**
 * @module ConfigurableSecretEngineDetails
 * The `ConfigurableSecretEngineDetails` is used by configurable secret engines to show either a prompt, error
 * or configuration details depending on the response from the engines specific config endpoint (ex: aws -> aws/root-config vs ssh: ssh/ca-config).
 *
 * @example ```js
 *   <ConfigurableSecretEngineDetails @model={{this.model}} />```
 *
 * @param {object} model - The secret-engine model to be configured.
 *
 */

interface Args {
  model: SecretEngineModel;
}

export default class ConfigurableSecretEngineDetails extends Component<Args> {
  @service declare readonly store: Store;
  @tracked configModel = null;
  @tracked configError: string | null = null;

  constructor(owner: unknown, args: Args) {
    super(owner, args);
    const { model } = this.args;
    if (!model) {
      this.configError =
        'We are unable to access the mount information for this engine. Ask you administrator if you think you should have access to this secret engine.';
      return;
    }
    // Fetch the config for the engine. Will eventually include GCP and Azure.
    if (model.type === 'aws') {
      this.fetchAwsRootConfig(model.id);
    }
    if (model.type === 'ssh') {
      this.fetchSshCaConfig(model.id);
    }
  }

  async fetchAwsRootConfig(backend: string) {
    try {
      this.configModel = await this.store.queryRecord('aws/root-config', { backend });
    } catch (e: AdapterError) {
      // If it's a 404 then the user has not configured the backend yet and we want to show the prompt instead
      if (e.httpStatus !== 404) {
        this.configError = errorMessage(e);
      }
      return;
    }
  }

  async fetchSshCaConfig(backend: string) {
    try {
      this.configModel = await this.store.queryRecord('ssh/ca-config', { backend });
    } catch (e: AdapterError) {
      // The SSH Api does not return a 404 not found but a 400 error after first mounting the engine with the
      // message that keys have not been configured yet.
      // To show a prompt instead of an error when first configuring the backend, we need to catch that specific 400 error and continue to set a prompt message instead.
      if (e.httpStatus !== 404 && errorMessage(e) !== `keys haven't been configured yet`) {
        this.configError = errorMessage(e);
      }
      return;
    }
  }

  get typeDisplay() {
    // Will eventually handle GCP and Azure.
    // Did not use capitalize helper because some are all caps and some only title case.
    const { type } = this.args.model;
    switch (type) {
      case 'aws':
        return 'AWS';
      case 'ssh':
        return 'SSH';
      default:
        return type;
    }
  }
}
