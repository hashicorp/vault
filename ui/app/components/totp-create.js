/**
 * Copyright (c) HashiCorp, Inc.
 * SPDX-License-Identifier: BUSL-1.1
 */

/**
 * @module TotpCreate
 * TotpCreate component displays the form for creating a new account
 *
 * @example
 * ```js
 * <TotpCreate
 *  @model={{model}}
 *  @modelForData={{@modelForData}}
 *  @error={{this.error}}
 *  @buttonDisabled={{this.saving}}
 * />
 * ```
 * @param {object} model - the route model
 * @param {object} modelForData - a class that helps track secret data, defined in secret-edit
 * @param {string} error - error message to be displayed
 * @param {boolean} buttonDisabled - if true, disables the submit button on the create/update form
 */

import Component from '@glimmer/component';
import ControlGroupError from 'vault/lib/control-group-error';
import Ember from 'ember';
import keys from 'core/utils/key-codes';
import { action, set } from '@ember/object';
import { service } from '@ember/service';
import { tracked } from '@glimmer/tracking';
import { isBlank } from '@ember/utils';
import { task, waitForEvent } from 'ember-concurrency';

const LIST_ROUTE = 'vault.cluster.secrets.backend.list';
const LIST_ROOT_ROUTE = 'vault.cluster.secrets.backend.list-root';
const SHOW_ROUTE = 'vault.cluster.secrets.backend.show';

export default class TotpCreate extends Component {
  @tracked error = null;
  @tracked nameWhiteSpaceWarning = false;
  @tracked validationErrorCount = 0;
  @tracked validationMessages = {
    name: '',
    url: '',
  };

  @service controlGroup;
  @service flashMessages;
  @service router;
  @service store;

  checkValidation(name, value) {
    if (name === 'name') {
      // check for whitespace
      this.nameHasWhiteSpace(value);

      if (!value) {
        set(this.validationMessages, name, `${name} can't be blank.`);
      } else if (value.includes('/')) {
        set(this.validationMessages, name, `${name} can't contain '/'.`);
      } else {
        set(this.validationMessages, name, '');
      }
    } else if (name === 'url') {
      if (!value) {
        set(this.validationMessages, name, `${name} can't be blank.`);
      } else {
        set(this.validationMessages, name, '');
      }
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

  nameHasWhiteSpace(value) {
    const validation = new RegExp('\\s', 'g'); // search for whitespace
    this.nameWhiteSpaceWarning = validation.test(value);
  }

  // successCallback is called in the context of the component
  persistKey(successCallback) {
    const secret = this.args.model;
    const secretData = this.args.modelForData;

    let key = secretData?.name;

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
  createKey(event) {
    event.preventDefault();
    if (isBlank(this.args.modelForData.name || this.args.modelForData.id)) {
      this.checkValidation('name', '');
      return;
    } else if (isBlank(this.args.modelForData.url)) {
      this.checkValidation('url', '');
      return;
    }

    const secretName = this.args.modelForData.name;
    this.persistKey(() => {
      // Show flash message in case there's a control group on read
      this.flashMessages.success(`Account ${secretName} added successfully.`);
      this.transitionToRoute(SHOW_ROUTE, secretName);
    });
  }

  @action
  updateValidationErrorCount(errorCount) {
    this.validationErrorCount = errorCount;
  }
}
