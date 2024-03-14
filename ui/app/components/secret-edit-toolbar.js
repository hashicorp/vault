/**
 * Copyright (c) HashiCorp, Inc.
 * SPDX-License-Identifier: BUSL-1.1
 */

/**
 * @module SecretEditToolbar
 * SecretEditToolbar component is the toolbar component displaying the JSON toggle and the actions like delete in the show mode.
 *
 * @example
 * ```js
 * <SecretEditToolbar
 * @mode={{mode}}
 * @model={{this.model}}
 * @isWriteWithoutRead={{isWriteWithoutRead}}
 * @secretDataIsAdvanced={{secretDataIsAdvanced}}
 * @showAdvancedMode={{showAdvancedMode}}
 * @modelForData={{this.modelForData}}
 * @canUpdateSecret={{canUpdateSecret}}
 * @editActions={{hash
    toggleAdvanced=(action "toggleAdvanced")
    refresh=(action "refresh")
  }}
 * />
 * ```

 * @param {string} mode - show, create, edit. The view.
 * @param {object} model - the model passed from the parent secret-edit
 * @param {boolean} isWriteWithoutRead - boolean describing permissions
 * @param {boolean} secretDataIsAdvanced - used to determine if show JSON toggle
 * @param {boolean} showAdvancedMode - used for JSON toggle
 * @param {object} modelForData - a modified version of the model with secret data
 * @param {boolean} canUpdateSecret - permissions to hide/show edit secret button.
 * @param {object} editActions - actions passed from parent to child
 */
/* eslint ember/no-computed-properties-in-native-classes: 'warn' */
import Component from '@glimmer/component';
import { action } from '@ember/object';
import { service } from '@ember/service';
import { tracked } from '@glimmer/tracking';
import { task } from 'ember-concurrency';
import { waitFor } from '@ember/test-waiters';

export default class SecretEditToolbar extends Component {
  @service store;
  @service router;
  @service flashMessages;

  @tracked wrappedData = null;

  @action
  clearWrappedData() {
    this.wrappedData = null;
  }

  @action
  handleDelete() {
    this.args.model.destroyRecord().then(() => {
      this.router.transitionTo('vault.cluster.secrets.backend.list-root');
    });
  }

  @task
  @waitFor
  *wrapSecret() {
    const { id } = this.args.modelForData;
    const { backend } = this.args.model;
    const wrapTTL = { wrapTTL: 1800 };

    try {
      const resp = yield this.store.adapterFor('secret').queryRecord(null, null, { backend, id, ...wrapTTL });
      this.wrappedData = resp.wrap_info.token;
      this.flashMessages.success('Secret successfully wrapped!');
    } catch (e) {
      this.flashMessages.danger('Could not wrap secret.');
    }
  }
}
