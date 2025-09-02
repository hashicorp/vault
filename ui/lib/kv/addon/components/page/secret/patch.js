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
 * @module KvSecretPatch
 * @description
 * This page template provides two methods for submitting patch data to update a KV v2 secret.
 * Either using a key/value form KvPatch::Editor::Form or the json editor via KvPatch::JsonForm
 *
 * <Page::Secret::Patch
 *  @backend="my-kv-engine"
 *  @breadcrumbs={{this.breadcrumbs}
 *  @metadata={{this.model.metadata}}
 *  @path="my-secret"
 *  @subkeys={{this.subkeys}
 *  @subkeysMeta={{this.subkeysMeta}
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
  @service('app-router') router;
  @service store;

  @tracked controlGroupError;
  @tracked errorMessage;
  @tracked invalidFormAlert;
  @tracked patchMethod = 'UI';

  @action
  selectPatchMethod(event) {
    this.patchMethod = event.target.value;
  }

  @task
  @waitFor
  *save(patchData) {
    const isEmpty = this.isEmpty(patchData);
    if (isEmpty) {
      this.flashMessages.info(`No changes to submit. No updates made to "${this.args.path}".`);
      return this.onCancel();
    }

    const { backend, path, metadata, subkeysMeta } = this.args;
    // if no metadata permission, use subkey metadata as backup
    const version = metadata?.currentVersion || subkeysMeta?.version;
    const adapter = this.store.adapterFor('kv/data');
    try {
      yield adapter.patchSecret(backend, path, patchData, version);
      this.flashMessages.success(`Successfully patched new version of ${path}.`);
      this.router.transitionTo('vault.cluster.secrets.backend.kv.secret.index');
    } catch (error) {
      if (error.message === 'Control Group encountered') {
        this.controlGroup.saveTokenFromError(error);
        this.controlGroupError = this.controlGroup.logFromError(error);
        return;
      }
      this.errorMessage = errorMessage(error);
      this.invalidFormAlert = 'There was an error submitting this form.';
    }
  }

  @action
  onCancel() {
    this.router.transitionTo('vault.cluster.secrets.backend.kv.secret.index');
  }

  isEmpty(object) {
    const emptyKeys = Object.keys(object).every((k) => k === '');
    const emptyValues = Object.values(object).every((v) => v === '');
    return emptyKeys && emptyValues;
  }
}
