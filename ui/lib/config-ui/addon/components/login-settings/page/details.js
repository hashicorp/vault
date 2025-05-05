/**
 * Copyright (c) HashiCorp, Inc.
 * SPDX-License-Identifier: BUSL-1.1
 */

import Component from '@glimmer/component';
import { service } from '@ember/service';

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
}
