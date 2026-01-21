/**
 * Copyright IBM Corp. 2016, 2025
 * SPDX-License-Identifier: BUSL-1.1
 */

import Component from '@glimmer/component';
import { service } from '@ember/service';

import type SecretMountPath from 'vault/services/secret-mount-path';

export default class HeaderScopeComponent extends Component {
  @service declare readonly secretMountPath: SecretMountPath;
}
