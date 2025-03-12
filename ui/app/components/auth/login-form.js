/**
 * Copyright (c) HashiCorp, Inc.
 * SPDX-License-Identifier: BUSL-1.1
 */

import Component from '@glimmer/component';
import { service } from '@ember/service';
import { tracked } from '@glimmer/tracking';
import { action } from '@ember/object';

import { task } from 'ember-concurrency';
import { waitFor } from '@ember/test-waiters';
import { allSupportedAuthBackends, supportedTypes } from 'vault/helpers/supported-auth-backends';

/**
 * @module Auth::LoginForm 2.0
 * */

export default class AuthLoginForm extends Component {
  @service flags;
  @service store;
  @service version;

  // form display logic
  @tracked authTabs = null; // listing visibility or backup types
  @tracked selectedTabIndex = 0;
  @tracked showOtherMethods = true;

  // auth login variables
  @tracked selectedAuthMethod = null;
  @tracked errorMessage = null;

  displayName = (type) => allSupportedAuthBackends().find((t) => t.type === type).typeDisplay;

  constructor() {
    super(...arguments);
    // todo fetch auth customization config here
    this.fetchMounts.perform();
  }

  get availableMethodTypes() {
    return supportedTypes(this.version.isEnterprise);
  }

  get namespaceInput() {
    const namespaceQP = this.args.namespace;
    if (this.flags.hvdManagedNamespaceRoot) {
      // When managed, the user isn't allowed to edit the prefix `admin/` for their nested namespace
      const split = namespaceQP.split('/');
      if (split.length > 1) {
        split.shift();
        return `/${split.join('/')}`;
      }
      return '';
    }
    return namespaceQP;
  }

  @action
  handleAuthSelect(element, event, idx) {
    if (element === 'tab') {
      this.selectedAuthMethod = Object.keys(this.authTabs)[idx];
      this.selectedTabIndex = idx;
    } else {
      this.selectedAuthMethod = event.target.value;
    }
  }

  @action
  handleError(message) {
    this.errorMessage = message;
  }

  @action
  onBackClick() {
    this.showOtherMethods = false;

    if (this.authTabs) {
      // go back to tabs
      // reset selected auth method to first tab
      this.handleAuthSelect('tab', null, 0);
    }
    // if default method exists, this takes us to single default method
    // do something about that here
  }

  @action
  onShowOtherClick() {
    this.showOtherMethods = true;

    if (this.authTabs) {
      // all methods render, reset dropdown
      this.selectedAuthMethod = 'token';
    }
    // if default method set, this takes us to "backups"
    // do something about that here
  }

  // this will be SO MUCH NICER with the auth updates that remove ember data
  fetchMounts = task(
    waitFor(async () => {
      try {
        const unauthMounts = await this.store.findAll('auth-method', {
          adapterOptions: {
            unauthenticated: true,
          },
        });

        if (unauthMounts.length !== 0) {
          this.authTabs = unauthMounts.reduce((obj, m) => {
            // serialize the ember data model into a regular ol' object
            const mountData = m.serialize();
            const methodType = mountData.type;
            if (!Object.keys(obj).includes(methodType)) {
              // create a new empty array for that type
              obj[methodType] = [];
            }
            // push mount data into corresponding type's array
            obj[methodType].push(mountData);
            return obj;
          }, {});

          // set first tab as selected method
          this.selectedAuthMethod = Object.keys(this.authTabs)[0];
          // hide other methods to prioritize tabs (visible mounts)
          this.showOtherMethods = false;
        }
      } catch (e) {
        // if for some reason there's an error fetching mounts, swallow and just show standard form
      }
    })
  );
}
