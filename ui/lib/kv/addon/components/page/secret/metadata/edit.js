/**
 * Copyright (c) HashiCorp, Inc.
 * SPDX-License-Identifier: BUSL-1.1
 */

import Component from '@glimmer/component';
import { service } from '@ember/service';
import { task } from 'ember-concurrency';
import { tracked } from '@glimmer/tracking';

/**
 * @module KvSecretMetadataEdit
 * This component renders the view for editing a kv secret's metadata.
 * While secret data and metadata are created on the same view, they are edited on different views/routes.
 *
 * @param {Form} form - kv form
 * @param {string} backend - mount path of the kv secret engine
 * @param {array} breadcrumbs - Breadcrumbs as an array of objects that contain label, route, and modelId. They are updated via the util kv-breadcrumbs to handle dynamic *pathToSecret on the list-directory route.
 * @param {object} capabilities - capabilities for data, metadata, subkeys, delete and undelete paths
 * @callback onCancel - Callback triggered when cancel button is clicked that transitions to the metadata details route.
 * @callback onSave - Callback triggered on save success that transitions to the metadata details route.
 */

export default class KvSecretMetadataEditComponent extends Component {
  @service flashMessages;
  @service api;

  @tracked errorBanner = '';
  @tracked invalidFormAlert = '';
  @tracked modelValidations = null;

  @task
  *save(event) {
    event.preventDefault();
    try {
      const { isValid, state, invalidFormMessage, data } = this.args.form.toJSON();
      this.modelValidations = isValid ? null : state;
      this.invalidFormAlert = invalidFormMessage;

      if (isValid) {
        const { path, ...metadata } = data;
        yield this.api.secrets.kvV2WriteMetadata(path, this.args.backend, metadata);
        this.flashMessages.success(`Successfully updated ${path}'s metadata.`);
        this.args.onSave();
      }
    } catch (error) {
      const { message } = yield this.api.parseError(error);
      this.errorBanner = message;
      this.invalidFormAlert = 'There was an error submitting this form.';
    }
  }
}
