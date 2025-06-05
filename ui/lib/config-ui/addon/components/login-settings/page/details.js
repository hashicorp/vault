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
 * @module Page::LoginSettingsRuleDetails
 * Page::LoginSettingsRuleDetails component is used to display rule information.
 * Shows detailed data, (eg which namespace it applies to, auth type used etc.) on a selected login rule from the custom login settings list.
 *
 * @example
 * ```js
 * <Page::LoginSettingsRuleDetails @rule={{this.rule}}  />
 * ```
 * @param {object} rule - holds login rule data, { backup_auth_types: string[] eg. ['token'], default_auth_type: string "oidc", disable_inheritance: boolean,
 * name: string "Login rule 1", namespace: string eg "admin/"}
 * */

export default class LoginSettingsRuleDetails extends Component {
  @service capabilities;
  @service flashMessages;
  @service('app-router') router;
  @service api;

  @tracked showConfirmModal = false;

  displayFields = {
    defaultAuthType: 'Default method',
    backupAuthTypes: 'Backup methods',
    disableInheritance: 'Inheritance enabled',
    namespacePath: 'Namespace the rule applies to',
  };

  displayValue = (key) => {
    const value = this.args.rule[key];
    // flip boolean for disable inheritance so template reads "Inheritance enabled: Yes/No"
    return key === 'disableInheritance' ? !value : value;
  };

  @action
  async onDelete() {
    const { rule } = this.args;

    try {
      await this.api.sys.uiLoginDefaultAuthDeleteConfiguration(rule.name);

      this.flashMessages.success(`Successfully deleted rule ${rule.name}.`);

      this.router.transitionTo('vault.cluster.config-ui.login-settings.index');
    } catch (error) {
      const message = errorMessage(error, 'Error deleting rule. Please try again.');
      this.flashMessages.danger(message);
    }
  }
}
