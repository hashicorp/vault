/**
 * Copyright (c) HashiCorp, Inc.
 * SPDX-License-Identifier: BUSL-1.1
 */

import { service } from '@ember/service';
import Component from '@glimmer/component';
import { tracked } from '@glimmer/tracking';
import { allEngines } from 'vault/helpers/mountable-secret-engines';

import type Store from '@ember-data/store';
import type SecretEngineModel from 'vault/models/secret-engine';
import type AwsRootConfig from 'vault/models/aws/root-config';
import type AwsLeaseConfig from 'vault/models/aws/lease-config';
import type SshCaConfig from 'vault/models/ssh/ca-config';
import type Model from '@ember-data/model';
import AdapterError from '@ember-data/adapter';

/**
 * @module ConfigurationDetails
 * `ConfigurationDetails` is used by configurable secret engines (AWS, SSH) to show either an API error, configuration details, or a prompt to configure the engine. Which of these is shown is determined by the engine type and whether the user has configured the engine yet.
 *
 * @example
 * ```js
 * <SecretEngine::ConfigurationDetails @model={{this.model}} />
 * ```
 *
 * @param {object} model - The secret-engine model to be configured.
 */

interface Args {
  model: SecretEngineModel | null;
  configModel: AwsLeaseConfig | AwsRootConfig | SshCaConfig | null;
}

interface ConfigError {
  httpStatus: number | null;
  message: string | null;
  errors: object | null;
}

export default class ConfigurationDetails extends Component<Args> {
  @service declare readonly store: Store;
  @tracked configError: ConfigError | null = null;
  @tracked configModel: Model | null = null;

  constructor(owner: unknown, args: Args) {
    super(owner, args);
    const { model } = this.args;
    // Should not be able to get here without a model, but in case an upstream change allows it, handle the error higher up.
    if (!model) {
      return;
    }
    // check if the configModel is an Error and what kind.
    const { configModel } = this.args;
    if (!configModel) return;

    if (configModel instanceof AdapterError) {
      if (configModel.httpStatus === 404 && configModel.errors[0] !== `keys haven't been configured yet`) {
        // If the error is 404, the user has not configured the engine yet.
        // SSH engines return a 400 error if the keys have not been configured yet. So check specifically for that message before assigning the error.
        this.configError = configModel;
        return;
      }
      return;
    }
    this.configModel = configModel;
  }

  get typeDisplay() {
    if (!this.args.model) return;
    const { type } = this.args.model;
    return allEngines().find((engine) => engine.type === type)?.displayName;
  }
}
