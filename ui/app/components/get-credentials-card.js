/**
 * Copyright (c) HashiCorp, Inc.
 * SPDX-License-Identifier: BUSL-1.1
 */

/**
 * @module GetCredentialsCard
 * Card-like component that display a title, and SearchSelect component that sends you to a credentials route for the selected item.
 * They are designed to be used in containers that act as flexbox or css grid containers. At this time, only used in the database
 * overview page to select a role and generate credentials
 *
 * @example
 * ```js
 * <GetCredentialsCard
 * @title="Get Credentials"
 * @searchLabel="Role to use"
 * @models={{array 'database/roles'}}
 * @backend={{model.backend}}
 * />
 * ```
 * @param {string} title - The title displays the card title
 * @param {string} searchLabel - The text above the searchSelect component
 * @param {array} models - An array of model types to fetch from the API. Passed through to SearchSelect component
 * @param {string} [placeholder] - Input placeholder text (default for SearchSelect is 'Search', none for InputSearch)
 * @param {string} backend - Passed to SearchSelect query method to fetch dropdown options
 */

import Component from '@glimmer/component';
import { service } from '@ember/service';
import { action } from '@ember/object';
import { tracked } from '@glimmer/tracking';

export default class GetCredentialsCard extends Component {
  @service router;
  @tracked role = '';

  @action
  transitionToCredential(evt) {
    evt.preventDefault();
    const role = this.role;
    if (role) {
      this.router.transitionTo('vault.cluster.secrets.backend.credentials', role);
    }
  }

  @action
  handleInput(value) {
    // if it comes in from the fallback component then the value is a string otherwise it's an array
    if (Array.isArray(value)) {
      this.role = value[0];
    } else {
      this.role = value;
    }
  }
}
