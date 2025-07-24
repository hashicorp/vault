/**
 * Copyright (c) HashiCorp, Inc.
 * SPDX-License-Identifier: BUSL-1.1
 */

import Component from '@glimmer/component';
import ControlGroupError from 'vault/lib/control-group-error';
import Ember from 'ember';
import keys from 'core/utils/keys';
import { action, set } from '@ember/object';
import { service } from '@ember/service';
import { tracked } from '@glimmer/tracking';
import { isBlank, isNone } from '@ember/utils';
import { task, waitForEvent } from 'ember-concurrency';
import { WHITESPACE_WARNING, noWhitespace } from 'vault/utils/forms/validators';

/**
 * @module SecretCreateOrUpdate
 * Component for creating or updating a secret.
 *
 * @example
 * ```hbs
 * <SecretCreateOrUpdate
 *   @mode="create"
 *   @model={{model}}
 *   @showAdvancedMode=true
 *   @modelForData={{@modelForData}}
 *   @secretData={{@secretData}}
 *   @buttonDisabled={{this.saving}}
 * />
 * ```
 *
 * @param {string} mode - "create", "edit", "show" to control the view.
 * @param {object} model - The route model.
 * @param {boolean} showAdvancedMode - Show JSON editor if true.
 * @param {object} modelForData - Helper class for managing secret data.
 * @param {object} secretData - KV data for the secret.
 * @param {boolean} buttonDisabled - Disables submit button when true.
 */

const LIST_ROUTE = 'vault.cluster.secrets.backend.list';
const LIST_ROOT_ROUTE = 'vault.cluster.secrets.backend.list-root';
const SHOW_ROUTE = 'vault.cluster.secrets.backend.show';

export default class SecretCreateOrUpdate extends Component {
  @tracked editorString = null;
  @tracked error = null;
  @tracked secretPaths = null;
  @tracked pathWhiteSpaceWarning = false;
  @tracked validationErrorCount = 0;
  @tracked validationMessages = null;

  @service controlGroup;
  @service flashMessages;
  @service router;
  @service store;

  whitespaceWarning = WHITESPACE_WARNING('path');

  @action
  setup(elem, [secretData, mode]) {
    this.editorString = secretData.toJSONString();
    this.validationMessages = { path: '' };

    if (Ember.testing) {
      this.secretPaths = ['beep', 'bop', 'boop'];
    }
    this.checkRows();
    if (mode === 'edit') this.addRow();
  }

  checkRows() {
    if (this.args.secretData.length === 0) {
      this.addRow();
    }
  }

  checkValidation(name, value) {
    if (name === 'path') {
      // Warning if whitespace exists
      this.pathWhiteSpaceWarning = !noWhitespace(value);

      if (!value) {
        set(this.validationMessages, name, `${name} can't be blank.`);
      } else if (!noWhitespace(value)) {
        set(this.validationMessages, name, this.whitespaceWarning);
      } else {
        set(this.validationMessages, name, '');
      }
    }
    this.validationErrorCount = Object.values(this.validationMessages).filter(Boolean).length;
  }

  onEscape(e) {
    const isEscKeyPressed = keys.ESC.includes(e.key);
    if (isEscKeyPressed || this.args.mode !== 'show') return;

    const parentKey = this.args.model.parentKey;
    this.transitionToRoute(parentKey ? LIST_ROUTE : LIST_ROOT_ROUTE, parentKey);
  }

  transitionToRoute() {
    return this.router.transitionTo(...arguments);
  }

  persistKey(successCallback) {
    const secret = this.args.model;
    const secretData = this.args.modelForData;
    let key = secretData?.path || secret.id;

    if (key.startsWith('/')) {
      key = key.replace(/^\/+/g, '');
      secretData.set(secretData.pathAttr, key);
    }

    return secretData
      .save()
      .then(() => {
        if (!secretData.isError) {
          this.saveComplete(successCallback, key);
        }
      })
      .catch((error) => {
        if (error instanceof ControlGroupError) {
          const errorMessage = this.controlGroup.logFromError(error);
          this.error = errorMessage.content;
          this.controlGroup.saveTokenFromError(error);
        }
        throw error;
      });
  }

  saveComplete(callback, key) {
    callback(key);
  }

  @(task(function* (name, value) {
    this.checkValidation(name, value);
    while (true) {
      const event = yield waitForEvent(document.body, 'keyup');
      this.onEscape(event);
    }
  })
    .on('didInsertElement')
    .cancelOn('willDestroyElement'))
  waitForKeyUp;

  @action
  addRow() {
    const data = this.args.secretData;
    if (isNone(data.find((d) => d.name === ''))) {
      data.pushObject({ name: '', value: '' });
      this.handleChange();
    }
    this.checkRows();
  }

  @action
  editorUpdated(val) {
    try {
      this.args.secretData.fromJSONString(val);
      set(this.args.modelForData, 'secretData', this.args.secretData.toJSON());
    } catch (e) {
      this.error = e.message;
    }
    this.editorString = val;
  }

  @action
  createOrUpdateKey(type, event) {
    event.preventDefault();
    if (type === 'create' && isBlank(this.args.modelForData.path || this.args.modelForData.id)) {
      this.checkValidation('path', '');
      return;
    }

    const secretPath = type === 'create' ? this.args.modelForData.path : this.args.model.id;
    this.persistKey(() => {
      this.flashMessages.success(
        `Secret ${secretPath} ${type === 'create' ? 'created' : 'updated'} successfully.`
      );
      this.transitionToRoute(SHOW_ROUTE, secretPath);
    });
  }

  @action
  deleteRow(name) {
    const data = this.args.secretData;
    const item = data.find((d) => d.name === name);
    if (isBlank(item.name)) return;
    data.removeObject(item);
    this.checkRows();
    this.handleChange();
  }

  @action
  formatJSON() {
    this.editorString = this.args.secretData.toJSONString(true);
  }

  @action
  handleMaskedInputChange(secret, index, value) {
    const row = { ...secret, value };
    set(this.args.secretData, index, row);
    this.handleChange();
  }

  @action
  handleChange() {
    this.editorString = this.args.secretData.toJSONString(true);
    set(this.args.modelForData, 'secretData', this.args.secretData.toJSON());
  }

  @action
  updateValidationErrorCount(errorCount) {
    this.validationErrorCount = errorCount;
  }
}
