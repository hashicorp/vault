/**
 * Copyright (c) HashiCorp, Inc.
 * SPDX-License-Identifier: BUSL-1.1
 */

import Component from '@glimmer/component';
import { service } from '@ember/service';

/**
 * @module Page::LoginSettingsList
 * Page::LoginSettingsList components are used to display list of rules.
 * @example
 * ```js
 * <Page::LoginSettingsList @loginRules={{this.rules}}  />
 * ```
 * @param {array} loginRules - array of rule objects
 */

export default class LoginSettingsList extends Component {
  @service capabilities;
}
