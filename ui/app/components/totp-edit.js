/**
 * Copyright (c) HashiCorp, Inc.
 * SPDX-License-Identifier: BUSL-1.1
 */

import Component from '@glimmer/component';
import { action } from '@ember/object';
import { tracked } from '@glimmer/tracking';
import { service } from '@ember/service';
import { task } from 'ember-concurrency';
import { waitFor } from '@ember/test-waiters';
import errorMessage from 'vault/utils/error-message';

/**
 * @module TotpEdit
 * `TotpEdit` is a component that allows you to create, view or delete a TOTP key.
 * When creating a key if `generate` and `exported` are true then after a successful save the UI renders a QR code for the generated key.
 * @example
 *   <TotpEdit @model={{this.model}} @mode={{this.mode}} />
 *
 * @param {object} model - The totp key ember data model.
 * @param {string} mode - The mode to render. Either 'create' or 'show'.
 */
const LIST_ROOT_ROUTE = 'vault.cluster.secrets.backend.list-root';
const SHOW_ROUTE = 'vault.cluster.secrets.backend.show';

export default class TotpEdit extends Component {
  @service router;
  @service flashMessages;

  @tracked hasGenerated = false;
  @tracked invalidFormAlert = '';
  @tracked modelValidations;

  successCallback;

  get defaultKeyFormFields() {
    const shared = ['name', 'generate', 'issuer', 'accountName'];
    const generated = [...shared, 'exported'];
    const nonGenerated = [...shared, 'url', 'key'];
    return this.args.model.generate ? generated : nonGenerated;
  }

  get groups() {
    const { generate } = this.args.model;

    const groups = {
      'TOTP Code Options': ['algorithm', 'digits', 'period'],
    };

    if (generate) {
      groups['Provider Options'] = ['keySize', 'skew', 'qrSize'];
    }

    return groups;
  }

  transitionToRoute() {
    this.router.transitionTo(...arguments);
  }

  @action
  reset() {
    const { name } = this.args.model;
    this.args.model.unloadRecord();
    this.transitionToRoute(SHOW_ROUTE, name);
  }

  @action
  async deleteKey() {
    try {
      const { id } = this.args.model;
      await this.args.model.destroyRecord();
      this.transitionToRoute(LIST_ROOT_ROUTE);
      this.flashMessages.success(`${id} was successfully deleted.`);
    } catch (err) {
      this.flashMessages.danger(errorMessage(err));
    }
  }

  createKey = task(
    waitFor(async (event) => {
      event.preventDefault();
      const { isValid, state, invalidFormMessage } = this.args.model.validate();
      this.modelValidations = isValid ? null : state;
      this.invalidFormAlert = invalidFormMessage;

      if (!isValid) return;
      try {
        const allFields = [...this.defaultKeyFormFields, ...Object.values(this.groups).flat()];
        await this.args.model.save({
          adapterOptions: {
            keyFormFields: allFields,
          },
        });
        const { generate, exported } = this.args.model;

        if (generate && exported) {
          // stay in this template and show QR code returned from response
          this.hasGenerated = true;
        } else {
          // nothing is returned from response, transition to key details route
          this.transitionToRoute(SHOW_ROUTE, this.args.model.name);
        }
      } catch (err) {
        // err will display via model state
        return;
      }
      this.flashMessages.success('Successfully created key.');
    })
  );
}
