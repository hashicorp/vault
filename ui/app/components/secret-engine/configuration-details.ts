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

/**
 * @module ConfigurationDetails
 * `ConfigurationDetails` is used by configurable secret engines (AWS, SSH) to show either an API error, configuration details, or a prompt to configure the engine. Which of these is shown is determined by the engine type and whether the user has configured the engine yet.
 *
 * @example
 * ```js
 * <SecretEngine::ConfigurationDetails @model={{this.model.backedn}} @configModel={{this.model.configModel}} />
 * ```
 *
 * @param {object} model s- The secret-engine model as this.model.backend and the configuration models with their name.
 * @param {object} configModel - The config model to be configured.
 */

interface Args {
  models: [SecretEngineModel | AwsLeaseConfig | AwsRootConfig | SshCaConfig];
}

interface ConfigError {
  httpStatus: number | null;
  message: string | null;
  errors: object | null;
}
// ARG TODO work on naming use plurals for arrays.
export default class ConfigurationDetails extends Component<Args> {
  @service declare readonly store: Store;
  @tracked configError: [ConfigError] | [] = [];
  @tracked configModel: [Model] | [] = [];
  @tracked engineType: string | '' = '';

  constructor(owner: unknown, args: Args) {
    super(owner, args);
    const { models } = this.args;
    // Should not be able to get here without the secret-engine model, but in case an upstream change allows it, handle the error higher up.
    if (!models.backend) return;
    this.engineType = models.backend.type; // save this now because modifying the models object later will remove the backend model.
    delete models.backend;
    // for each configModel check if error, or not configured, or configured.
    Object.values(models).forEach((configModel) => {
      // ARG STOPPED HERE. Remember you just need to check if the configModel is an AdapterError and if so, assign it to the configError property.
      // Otherwise show one or two of the configModels.
      this.configModelAssignment(configModel);
    });
  }

  configModelAssignment(configModel: Model) {
    if (configModel.isAdapterError) {
      // Check for errors that indicate the engine has not been configured yet. If they haven't return nothing so that the form displays the configuration prompt.
      // Most engines return a 404 if they have not been configured, but SSH returns a 400 and a specific error message that we check for here.
      if (
        (this.engineType === 'ssh' &&
          configModel.httpStatus === 400 &&
          configModel.errors[0] === `keys haven't been configured yet`) ||
        configModel.httpStatus === 404
      ) {
        return;
      }
      this.configError.push(configModel);
      return;
    }
    this.configModel.push(configModel);
    return;
  }

  get typeDisplay() {
    if (!this.engineType) return;
    return allEngines().find((engine) => engine.type === this.engineType)?.displayName;
  }
}
