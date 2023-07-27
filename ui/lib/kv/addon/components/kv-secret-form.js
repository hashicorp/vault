/**
 * Copyright (c) HashiCorp, Inc.
 * SPDX-License-Identifier: MPL-2.0
 */

import Component from '@glimmer/component';
import { action } from '@ember/object';
import { tracked } from '@glimmer/tracking';
import { task } from 'ember-concurrency';
import { inject as service } from '@ember/service';
import KVObject from 'vault/lib/kv-object';

/**
 * @module KvSecretForm is used for creating a new secret or secret version (also considered 'editing')
 *
 * <KvSecretForm
 *  @secret={{@secret}}
 *  @onSave={{transition-to "vault.cluster.secrets.backend.kv.secret.details" @secret.path}}
 *  @onCancel={{transition-to "vault.cluster.secrets.backend.kv.list"}}
 * />
 *
 * @param {model} secret - Ember data model: 'kv/data'
 * @param {callback} onSave - callback (usually a transition) from parent to perform after the model is saved
 * @param {callback} onCancel - callback (usually a transition) from parent to perform when cancel button is clicked
 */

export default class KvSecretForm extends Component {
  @service flashMessages;
  @tracked showJsonView = false;
  @tracked errorMessage;
  @tracked modelValidations;
  @tracked invalidFormAlert;

  get emptyJson() {
    // if secretData is null, this specially formats a blank object and renders a nice initial state for the json editor
    return KVObject.create({ content: [{ name: '', value: '' }] }).toJSONString(true);
  }

  @action
  toggleJsonView() {
    this.showJsonView = !this.showJsonView;
  }

  @action
  handleJson(value, codemirror) {
    codemirror.performLint();
    const lintingErrors = codemirror.state.lint.marked.length > 0;
    if (!lintingErrors) {
      this.args.secret.secretData = JSON.parse(value);
    }
  }

  @action
  keyUpValidations() {
    // check warnings on key up
    const { state } = this.args.secret.validate();
    for (const attr in state) {
      state[attr].errors = [];
    }
    this.modelValidations = state;
  }

  @task
  *save(event) {
    event.preventDefault();
    try {
      const { isValid, state, invalidFormMessage } = this.args.secret.validate();
      this.modelValidations = isValid ? null : state;
      this.invalidFormAlert = invalidFormMessage;
      if (isValid) {
        const { path, isNew } = this.args.secret;
        yield this.args.secret.save();
        this.flashMessages.success(`Successfully created ${isNew ? '' : 'new version of'} secret ${path}`);
        this.args.onSave();
      }
    } catch (error) {
      const message = error.errors ? error.errors.join('. ') : error.message;
      this.errorMessage = message;
      this.invalidFormAlert = 'There was an error submitting this form.';
    }
  }

  @action
  cancel() {
    this.args.secret.rollbackAttributes();
    this.args.onCancel();
  }
}
