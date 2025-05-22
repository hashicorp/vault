/**
 * Copyright (c) HashiCorp, Inc.
 * SPDX-License-Identifier: BUSL-1.1
 */

import Component from '@glimmer/component';
import { tracked } from '@glimmer/tracking';
import { service } from '@ember/service';
import { action } from '@ember/object';
import errorMessage from 'vault/utils/error-message';

/**
 * @module Page::LoginSettingsList
 * Page::LoginSettingsList components are used to display list of rules.
 *
 * @example
 * <Page::LoginSettingsList @loginRules={{this.rules}}  />
 *
 * @param {array} loginRules - array of rule objects
 */

export default class LoginSettingsList extends Component {
  @service capabilities;
  @service flashMessages;
  @service('app-router') router;
  @tracked ruleToDelete = null; // set to the rule intended to delete
  @service api;

  @action
  async onDelete() {
    try {
      await this.api.sys.uiLoginDefaultAuthDeleteConfiguration(this.ruleToDelete.id);
      this.flashMessages.success(`Successfully deleted rule ${this.ruleToDelete.id}.`);

      this.router.transitionTo('vault.cluster.config-ui.login-settings');
    } catch (error) {
      const message = errorMessage(error, 'Error deleting rule. Please try again.');
      this.flashMessages.danger(message);
    }
  }
}
