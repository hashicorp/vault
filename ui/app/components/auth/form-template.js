/**
 * Copyright (c) HashiCorp, Inc.
 * SPDX-License-Identifier: BUSL-1.1
 */

import Component from '@glimmer/component';
import { service } from '@ember/service';
import { tracked } from '@glimmer/tracking';
import { action } from '@ember/object';
import { allSupportedAuthBackends, supportedAuthBackends } from 'vault/helpers/supported-auth-backends';
import { task } from 'ember-concurrency';

/**
 * @module Auth::FormTemplate
 *
 * @example
 *
 * @param {string} authMethodQueryParam - auth method type to login with, updated by selecting an auth method from the dropdown
 */

export default class AuthFormTemplate extends Component {
  @service version;
  @service store;
  @tracked signInWithAll = false; // TODO logic for base form (no customizations)
  @tracked logInWithOther = false; // TODO logic for base form (no customizations)

  @tracked authType = 'token';
  @tracked authTabs = null;

  displayName = (type) => this.allTypes.find((t) => t.type === type).typeDisplay;

  constructor() {
    super(...arguments);
    this.fetchMounts.perform();
  }

  get selectedAuth() {
    return this.args.selectedAuth;
  }

  get allTypes() {
    return this.version.isEnterprise ? allSupportedAuthBackends() : supportedAuthBackends();
  }

  get isDirectLink() {
    // TODO logic for whether URL contains ?with=
    return false;
  }

  fetchMounts = task(async () => {
    try {
      const unauthMounts = await this.store.findAll('auth-method', {
        adapterOptions: {
          unauthenticated: true,
        },
      });
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
        console;
        return obj;
      }, {});
    } catch (e) {
      // todo remove and swallow
      debugger;
    }
  });

  @action
  handleTypeChange(event) {
    console.log(event, 'hello?');
    this.authType = event.target.value;
  }

  @action
  handleError() {
    // do something
  }

  @action
  onTabClick(evt, idx) {
    const type = Object.keys(this.authTabs)[idx];
    this.args.onAuthChange(type);
  }

  @action
  handleInput(evt) {
    // For changing values in this backing class, not on form
    const { name, value } = evt.target;
    this[name] = value;

    if (this.args.onUpdate) {
      // Do parent side effects like update query params
      this.args.onUpdate(name, value);
    }
  }
}
