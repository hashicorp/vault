/**
 * Copyright IBM Corp. 2016, 2025
 * SPDX-License-Identifier: BUSL-1.1
 */

import Component from '@glimmer/component';
import { service } from '@ember/service';

import { action } from '@ember/object';
import { tracked } from '@glimmer/tracking';
import { task } from 'ember-concurrency';
import type KvForm from 'vault/app/forms/secrets/kv';
import type { Breadcrumb } from 'vault/vault/app-types';
import type FlashMessageService from 'vault/services/flash-messages';
import type RouterService from '@ember/routing/router-service';
import type ApiService from 'vault/services/api';
import SecretsEngineResource from 'vault/app/resources/secrets/engine';

interface Args {
  form: KvForm;
  backend: SecretsEngineResource;
  breadcrumbs: Array<Breadcrumb>;
}

/**
 * @module KvConfigurePageComponent
 * KvConfigurePageComponent is a component to show secrets mount and engine configuration data
 *
 * @param {object} form - config form data for mount and engine
 * @param {string} backend - The kv secrets engine data
 * @param {array} breadcrumbs - Breadcrumbs as an array of objects that contain label, route, and modelId. They are updated via the util kv-breadcrumbs to handle dynamic *pathToSecret on the list-directory route.
 */

export default class KvConfigurePageComponent extends Component<Args> {
  @service declare readonly flashMessages: FlashMessageService;
  @service('app-router') declare readonly router: RouterService;
  @service declare readonly api: ApiService;

  @tracked errorBanner = '';
  @tracked invalidFormAlert = '';
  @tracked modelValidations = null;

  @action
  navigateToConfiguration() {
    this.router.transitionTo(`vault.cluster.secrets.backend.kv.configuration`);
  }

  @task
  *save(event: Event | null) {
    event.preventDefault();
    try {
      const { isValid, state, invalidFormMessage, data } = this.args.form.toJSON();
      this.modelValidations = isValid ? null : state;
      this.invalidFormAlert = invalidFormMessage;

      if (isValid) {
        yield this.api.secrets.kvV2Configure(data.path, data);
        this.flashMessages.success(`Successfully updated ${data.path}'s configuration.`);
        this.navigateToConfiguration();
      }
    } catch (error) {
      const { message } = yield this.api.parseError(error);
      this.errorBanner = message;
      this.invalidFormAlert = 'There was an error submitting this form.';
    }
  }
}
