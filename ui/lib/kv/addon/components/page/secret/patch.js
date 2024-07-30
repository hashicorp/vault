/**
 * Copyright (c) HashiCorp, Inc.
 * SPDX-License-Identifier: BUSL-1.1
 */

import Component from '@glimmer/component';
import { service } from '@ember/service';
import { action } from '@ember/object';
import { tracked } from '@glimmer/tracking';
import { task } from 'ember-concurrency';
import { waitFor } from '@ember/test-waiters';
import errorMessage from 'vault/utils/error-message';

/**
 * @module KvSecretPatch is used for creating a new version of a secret
 *
 * <Page::Secret::Patch
 *  @path="my-secret"
 *  @backend="my-kv-engine"
 *  @metadata={{this.model.metadata}}
 *  @breadcrumbs={{this.breadcrumbs}
 *  @subkeys={{this.subkeys}
 * />
 *
 * @param {model} path - Secret path
 * @param {string} backend - Mount backend path
 * @param {model} metadata - Ember data model: 'kv/metadata'
 * @param {object} subkeys - subkeys (leaf keys with null values) of kv v2 secret
 * @param {object} subkeysMeta - metadata object returned from the /subkeys endpoint, contains: version, created_time, custom_metadata, deletion status and time
 * @param {array} breadcrumbs - breadcrumb objects to render in page header
 */

export default class KvSecretPatch extends Component {
  @service controlGroup;
  @service flashMessages;
  @service router;
  @service store;

  @tracked jsonObject;
  @tracked kvObject;
  @tracked lintingErrors;
  @tracked patchMethod = 'UI';
  _emptyJson = JSON.stringify({ '': '' }, null, 2);
  _emptySubkeys;

  constructor() {
    super(...arguments);
    this._emptySubkeys = Object.keys(this.args.subkeys).reduce((obj, key) => {
      obj[key] = '';
      return obj;
    }, {});
    this.resetForm();
  }

  get formData() {
    return this.patchMethod === 'UI' ? this.kvObject : this.jsonObject;
  }

  resetForm() {
    this.kvObject = this._emptySubkeys;
    this.jsonObject = this._emptyJson;
  }

  @action
  selectPatchMethod(event) {
    this.patchMethod = event.target.value;
    this.resetForm();
  }

  @action
  handleJson(value, codemirror) {
    codemirror.performLint();
    this.lintingErrors = codemirror.state.lint.marked.length > 0;
    if (!this.lintingErrors) {
      this.jsonObject = JSON.parse(value);
    }
  }

  @task
  @waitFor
  *submit(event) {
    event.preventDefault();
    const { backend, path, metadata, subkeysMeta } = this.args;
    const adapter = this.store.adapterFor('kv/data');
    // we can't guarantee the subkeys meta will be the latest version, so it's backup
    const version = metadata.currentVersion || subkeysMeta.version;
    const data = {
      options: { cas: version },
      data: this.formData,
    };
    try {
      yield adapter.patchSecret(backend, path, data);
      this.flashMessages.success(`Successfully patched new version of ${path}.`);
      this.router.transitionTo('vault.cluster.secrets.backend.kv.secret');
    } catch (error) {
      // TODO test...this is copy pasta'd from the edit page
      let message = errorMessage(error);
      if (error.message === 'Control Group encountered') {
        this.controlGroup.saveTokenFromError(error);
        const err = this.controlGroup.logFromError(error);
        message = err.content;
      }
      this.errorMessage = message;
      this.invalidFormAlert = 'There was an error submitting this form.';
    }
  }

  @action
  onCancel() {
    this.router.transitionTo('vault.cluster.secrets.backend.kv.secret');
  }
}
