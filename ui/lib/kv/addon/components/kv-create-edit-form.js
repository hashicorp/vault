/**
 * Copyright (c) HashiCorp, Inc.
 * SPDX-License-Identifier: BUSL-1.1
 */

import Component from '@glimmer/component';
import { service } from '@ember/service';
import { action } from '@ember/object';
import { tracked } from '@glimmer/tracking';
import { task } from 'ember-concurrency';
import { pathIsFromDirectory } from 'kv/utils/kv-breadcrumbs';
import { waitFor } from '@ember/test-waiters';

/**
 * @module KvCreateEditForm is used for creating and editing kv secret data and metadata, it hides/shows a json editor and renders validation errors for the json editor
 *
 * <KvCreateEditForm
 *  @form={{@form}}
 *  @path={{@path}}
 *  @backend={{@backend}}
 *  @showJson={{true}}
 *  @onChange={{@onChange}}
 * />
 *
 * @param {Form} form - kv form
 * @param {string} path - secret path
 * @param {string} backend - secret mount path
 * @param {boolean} showJson - boolean passed from parent to hide/show json editor
 * @param {function} onSecretDataChange - function passed from parent to handle secret data change side effects
 */

export default class KvCreateEditForm extends Component {
  @service api;
  @service controlGroup;
  @service flashMessages;
  @service('app-router') router;

  @tracked lintingErrors;
  @tracked modelValidations;
  @tracked invalidFormAlert;
  @tracked errorMessage;

  @action
  onJsonChange(value) {
    try {
      const json = JSON.parse(value);
      this.args.form.data.secretData = json;
      this.lintingErrors = false;
      this.args.onChange?.(json);
    } catch {
      this.lintingErrors = true;
    }
  }

  @action
  onKvObjectChange(value) {
    this.args.form.data.secretData = value;
    this.args.onChange?.(value);
  }

  @action
  pathValidations() {
    // check path attribute warnings on key up for new secrets
    const { state } = this.args.form.toJSON();
    if (state?.path?.warnings) {
      // only set model validations if warnings exist
      this.modelValidations = state;
    }
  }

  @action
  onCancel() {
    const { form, path } = this.args;
    if (form.isNew) {
      pathIsFromDirectory(path)
        ? this.router.transitionTo('vault.cluster.secrets.backend.kv.list-directory', path)
        : this.router.transitionTo('vault.cluster.secrets.backend.kv.list');
    } else {
      this.router.transitionTo('vault.cluster.secrets.backend.kv.secret.index');
    }
  }

  hasMetadata(metadata) {
    try {
      const { custom_metadata = {}, max_versions, cas_required, delete_version_after = '0s' } = metadata;
      return (
        Object.keys(custom_metadata).length || max_versions || cas_required || delete_version_after !== '0s'
      );
    } catch (e) {
      return false;
    }
  }

  save = task(
    waitFor(async (event) => {
      event.preventDefault();

      const { isValid, state, invalidFormMessage, data } = this.args.form.toJSON();
      this.modelValidations = isValid ? null : state;
      this.invalidFormAlert = invalidFormMessage;
      this.errorMessage = null;

      if (isValid) {
        const { path, secretData, options, ...metadata } = data;
        try {
          // try saving secret data first
          const payload = options ? { data: secretData, options } : { data: secretData };
          await this.api.secrets.kvV2Write(path, this.args.backend, payload);
          this.flashMessages.success(`Successfully saved secret data for: ${path}.`);

          // users must have permission to create secret data to create metadata in the UI
          // only attempt to save metadata if secret data saves successfully and metadata is added
          if (this.hasMetadata(metadata)) {
            try {
              await this.api.secrets.kvV2WriteMetadata(path, this.args.backend, metadata);
              this.flashMessages.success(`Successfully saved metadata.`);
            } catch (error) {
              const { message } = await this.api.parseError(error);
              this.flashMessages.danger(`Secret data was saved but metadata was not: ${message}`, {
                sticky: true,
              });
            }
          }
        } catch (error) {
          const { message, response } = await this.api.parseError(error);
          if (response.isControlGroupError) {
            this.controlGroup.saveTokenFromError(response);
            const err = this.controlGroup.logFromError(response);
            this.errorMessage = err.content;
          } else {
            this.errorMessage = message;
          }
        }

        // prevent transition if there are errors with secret data
        if (this.errorMessage) {
          this.invalidFormAlert = 'There was an error submitting this form.';
        } else {
          this.router.transitionTo('vault.cluster.secrets.backend.kv.secret.index', path);
        }
      }
    })
  );
}
