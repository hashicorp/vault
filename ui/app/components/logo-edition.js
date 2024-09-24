/**
 * Copyright (c) HashiCorp, Inc.
 * SPDX-License-Identifier: BUSL-1.1
 */

import { service } from '@ember/service';
import Component from '@ember/component';

/**
 * @module LogoEdition
 * LogoEdition shows the Vault logo with information about enterprise if applicable.
 *
 * @example
 * ```js
 * <LogoEdition />
 */

export default class LogoEdition extends Component {
  @service version;
}
