/**
 * Copyright (c) HashiCorp, Inc.
 * SPDX-License-Identifier: BUSL-1.1
 */

/**
 * @module SecretCreateOrUpdate
 * SecretCreateOrUpdate component displays either the form for creating a new secret or creating a new version of the secret
 *
 * @example
 * ```js
 * <SecretCreateOrUpdate
 *  @mode="create"
 *  @model={{model}}
 *  @showAdvancedMode=true
 *  @modelForData={{@modelForData}}
 *  @secretData={{@secretData}}
 *  @buttonDisabled={{this.saving}}
 * />
 * ```
 * @param {string} mode - create, edit, show determines what view to display
 * @param {object} model - the route model
 * @param {boolean} showAdvancedMode - whether or not to show the JSON editor
 * @param {object} modelForData - a class that helps track secret data, defined in secret-edit
 * @param {object} secretData - class that is created in secret-edit
 * @param {boolean} buttonDisabled - if true, disables the submit button on the create/update form
 */

import Component from '@glimmer/component';
import ControlGroupError from 'vault/lib/control-group-error';
import Ember from 'ember';
import keys from 'core/utils/key-codes';
import { action, set } from '@ember/object';
import { service } from '@ember/service';
import { tracked } from '@glimmer/tracking';
import { isBlank, isNone } from '@ember/utils';
import { task, waitForEvent } from 'ember-concurrency';

const LIST_ROUTE = 'vault.cluster.secrets.backend.list';
const LIST_ROOT_ROUTE = 'vault.cluster.secrets.backend.list-root';
const SHOW_ROUTE = 'vault.cluster.secrets.backend.show';

export default class SecretCreateOrUpdate extends Component {
  @tracked codemirrorString = null;
  @tracked error = null;
  @tracked secretPaths = null;
  @tracked pathWhiteSpaceWarning = false;
  @tracked validationErrorCount = 0;
  @tracked validationMessages = null;

  @service controlGroup;
  @service flashMessages;
  @service router;
  @service store;

  @action
  setup(elem, [secretData, mode]) {
    this.codemirrorString = secretData.toJSONString();
    this.validationMessages = {
      path: '',
    };
    // for validation, return array of path names already assigned
    if (Ember.testing) {
      this.secretPaths = ['beep', 'bop', 'boop'];
    }
    this.checkRows();

    if (mode === 'edit') {
      this.addRow();
    }
  }
  checkRows() {
    if (this.args.secretData.length === 0) {
      this.addRow();
    }
  }
  checkValidation(name, value) {
    if (name === 'path') {
      // check for whitespace
      this.pathHasWhiteSpace(value);
      !value
        ? set(this.validationMessages, name, `${name} can't be blank.`)
        : set(this.validationMessages, name, '');
    }
    const values = Object.values(this.validationMessages);
    this.validationErrorCount = values.filter(Boolean).length;
  }
  onEscape(e) {
    if (e.keyCode !== keys.ESC || this.args.mode !== 'show') {
      return;
    }
    const parentKey = this.args.model.parentKey;
    if (parentKey) {
      this.transitionToRoute(LIST_ROUTE, parentKey);
    } else {
      this.transitionToRoute(LIST_ROOT_ROUTE);
    }
  }
  pathHasWhiteSpace(value) {
    const validation = new RegExp('\\s', 'g'); // search for whitespace
    this.pathWhiteSpaceWarning = validation.test(value);
  }
  // successCallback is called in the context of the component
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
  transitionToRoute() {
    return this.router.transitionTo(...arguments);
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
    // fired off on init
    if (isNone(data.find((d) => d.name === ''))) {
      data.pushObject({ name: '', value: '' });
      this.handleChange();
    }
    this.checkRows();
  }
  @action
  codemirrorUpdated(val, codemirror) {
    this.error = null;
    codemirror.performLint();
    const noErrors = codemirror.state.lint.marked.length === 0;
    if (noErrors) {
      try {
        this.args.secretData.fromJSONString(val);
        set(this.args.modelForData, 'secretData', this.args.secretData.toJSON());
      } catch (e) {
        this.error = e.message;
      }
    }
    this.codemirrorString = val;
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
      // Show flash message in case there's a control group on read
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
    if (isBlank(item.name)) {
      return;
    }
    // secretData is a KVObject/ArrayProxy so removeObject is fine here
    data.removeObject(item);
    this.checkRows();
    this.handleChange();
  }
  @action
  formatJSON() {
    this.codemirrorString = this.args.secretData.toJSONString(true);
  }
  @action
  handleMaskedInputChange(secret, index, value) {
    const row = { ...secret, value };
    set(this.args.secretData, index, row);
    this.handleChange();
  }
  @action
  handleChange() {
    this.codemirrorString = this.args.secretData.toJSONString(true);
    set(this.args.modelForData, 'secretData', this.args.secretData.toJSON());
  }
  @action
  updateValidationErrorCount(errorCount) {
    this.validationErrorCount = errorCount;
  }
}
