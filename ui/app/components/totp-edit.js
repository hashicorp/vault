/**
 * Copyright (c) HashiCorp, Inc.
 * SPDX-License-Identifier: BUSL-1.1
 */

import Component from '@glimmer/component';
import { action } from '@ember/object';
import { tracked } from '@glimmer/tracking';
import { service } from '@ember/service';
import { isBlank } from '@ember/utils';

// TODO add jsdoc documentation
const LIST_ROOT_ROUTE = 'vault.cluster.secrets.backend.list-root';
const SHOW_ROUTE = 'vault.cluster.secrets.backend.show';

export default class TotpEdit extends Component {
  @service router;

  @tracked hasGenerated = false;
  successCallback;
  generatedDefaultFields = ['name', 'generate', 'issuer', 'accountName'];
  nonGeneratedDefaultFields = [...this.generatedDefaultFields, 'url', 'key'];

  get generatedFields() {
    return this.generatedDefaultFields;
  }

  get nonGeneratedFields() {
    return this.nonGeneratedDefaultFields;
  }

  get mode() {
    return this.args.mode || 'show';
  }

  get model() {
    return this.args.model;
  }

  persist(method, successCallback) {
    // TODO refactor this further
    return this.model[method]().then(() => {
      if (!this.model.isError) {
        if (this.model.backend === 'totp' && this.model.generate) {
          this.hasGenerated = true;
          this.successCallback = successCallback;
        } else {
          successCallback(this.model);
        }
      }
    });
  }

  transitionToRoute() {
    this.router.transitionTo(...arguments);
  }

  @action
  reset() {
    this.model.unloadRecord();
    this.successCallback(null);
  }

  @action
  delete() {
    this.persist('destroyRecord', () => {
      this.transitionToRoute(LIST_ROOT_ROUTE);
    });
  }

  @action
  create(event) {
    event.preventDefault();
    const modelId = this.model.name;
    // prevent from submitting if there's no key
    // maybe do something fancier later
    if (isBlank(modelId)) {
      return;
    }

    this.persist('save', () => {
      this.transitionToRoute(SHOW_ROUTE, modelId);
    });
  }
}
