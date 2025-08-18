/**
 * Copyright (c) HashiCorp, Inc.
 * SPDX-License-Identifier: BUSL-1.1
 */

import Component from '@glimmer/component';
import { task } from 'ember-concurrency';
import { service } from '@ember/service';
import { tracked } from '@glimmer/tracking';
import { ValidationMap } from 'vault/vault/app-types';

import type GeneralSettingsForm from 'vault/forms/secrets/general-settings';
import type Router from '@ember/routing/router';
import type FlashMessageService from 'vault/services/flash-messages';
import type ApiService from 'vault/services/api';
import type SecretsEngineResource from 'vault/resources/secrets/engine';

interface EngineData {
  displayName: string;
  type: string;
  glyph: string;
  mountCategory: String[];
}

/**
 * @module GeneralSettingsComponent is used to configure the SSH secret engine.
 *
 * @example
 * ```js
 * <Secrets:Page:GeneralSettings
 *    @secretsEngine={{@secretsEngine}}
 *    @engineData={{@engineData}}
 *    @generalSettingsForm={{@generalSettingsForm}}
 *  />
 * ```
 *
 * @param {string} secretsEngine - secrets engine resource
 * @param {string} engineData - hardcoded secrets engine metadata
 * @param {string} generalSettingsForm - general settings form
 */

interface Args {
  secretsEngine: SecretsEngineResource;
  engineData: EngineData;
  generalSettingsForm: GeneralSettingsForm;
}

export default class GeneralSettingsComponent extends Component<Args> {
  @service declare readonly router: Router;
  @service declare readonly api: ApiService;
  @service declare readonly flashMessages: FlashMessageService;

  @tracked errorMessage: string | null = null;
  @tracked invalidFormAlert: string | null = null;
  @tracked modelValidations: ValidationMap | null = null;

  saveGeneralSettings = task(async (event) => {
    event.preventDefault();
    // handle error state
    // make post request to tune endpoint
    console.log('save');
  });
}
