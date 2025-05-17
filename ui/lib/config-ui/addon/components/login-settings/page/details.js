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
  @service store;

  @tracked showConfirmModal = false;

  @action
  async onDelete() {
    const { rule } = this.args;

    try {
      const adapter = this.store.adapterFor('application');

      await adapter.ajax(`/v1/sys/config/ui/login/default-auth/${encodeURI(rule.name)}`, 'DELETE');
      this.flashMessages.success(`Successfully deleted rule ${rule.name}.`);

      this.router.transitionTo('vault.cluster.config-ui.login-settings.index');
    } catch (error) {
      const message = errorMessage(error, 'Error deleting rule. Please try again.');
      this.flashMessages.danger(message);
    }
  }
}
