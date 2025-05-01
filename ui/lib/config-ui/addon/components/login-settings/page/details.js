/**
 * Copyright (c) HashiCorp, Inc.
 * SPDX-License-Identifier: BUSL-1.1
 */

import Component from '@glimmer/component';
import { service } from '@ember/service';

/**
 * @module Page::LoginSettingsRuleDetails
 * Page::LoginSettingsRuleDetails components are used to display list of rules.
 * @example
 * ```js
 * <Page::LoginSettingsRuleDetails @rule={{this.rule}}  />
 * ```
 * @param {object} rule - rule object
 */

export default class LoginSettingsRuleDetails extends Component {
  @service capabilities;
}
