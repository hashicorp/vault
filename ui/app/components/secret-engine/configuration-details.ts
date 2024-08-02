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
 * <SecretEngine::ConfigurationDetails @configModels={{this.configModels}} />
 * ```
 * @param {string} typeDisplay - String of how we want to display the engine name (ex: SSH or Azure).
 * @param {string} id - Backend/path/name/id of the secret engine. Example: 'aws-123'.
 * @param {object} configModels - An object of config model(s).
 */

interface Args {
  models: Array<Model>;
}

interface ConfigError {
  httpStatus: number | null;
  message: string | null;
  errors: object | null;
}

export default class ConfigurationDetails extends Component<Args> {}
