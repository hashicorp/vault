/**
 * Copyright (c) HashiCorp, Inc.
 * SPDX-License-Identifier: MPL-2.0
 */

import Component from '@glimmer/component';
import { inject as service } from '@ember/service';
import { task } from 'ember-concurrency';
import { tracked } from '@glimmer/tracking';
import errorMessage from 'vault/utils/error-message';
import { action } from '@ember/object';

/**
 * @module KvSecretMetadataEdit
 * This component renders the view for editing a kv secret's metadata.
 * While secret data and metadata are created on the same view, they are edited on different views/routes.
 *
 * @param {array} metadata - The kv/metadata model. It is version agnostic.
 * @param {array} breadcrumbs - Breadcrumbs as an array of objects that contain label, route, and modelId. They are updated via the util kv-breadcrumbs to handle dynamic *pathToSecret on the list-directory route.
 * @callback onCancel - Callback triggered when cancel button is clicked that transitions to the metadata details route.
 * @callback onSave - Callback triggered on save success that transitions to the metadata details route.
 */

export default class KvSecretMetadataEditComponent extends Component {
  @service flashMessages;
  @tracked errorBanner = '';
  @tracked invalidFormAlert = '';
  @tracked modelValidations = null;

  @action
  cancel() {
    this.args.metadata.rollbackAttributes();
    this.args.onCancel();
  }

  @task
  *save(event) {
    event.preventDefault();
    try {
      const { isValid, state, invalidFormMessage } = this.args.metadata.validate();
      this.modelValidations = isValid ? null : state;
      this.invalidFormAlert = invalidFormMessage;
      if (isValid) {
        const { path } = this.args.metadata;
        yield this.args.metadata.save();
        this.flashMessages.success(`Successfully updated ${path}'s metadata.`);
        this.args.onSave();
      }
    } catch (error) {
      this.errorBanner = errorMessage(error);
      this.invalidFormAlert = 'There was an error submitting this form.';
    }
  }
}
