/**
 * Copyright IBM Corp. 2016, 2025
 * SPDX-License-Identifier: BUSL-1.1
 */

import Component from '@glimmer/component';
import { service } from '@ember/service';

import type SecretMountPath from 'vault/services/secret-mount-path';

interface Args {
  currentRoute: string;
  showPath?: boolean;
  scope?: string;
  role?: string;
}

export default class KmipBreadcrumbComponent extends Component<Args> {
  @service declare secretMountPath: SecretMountPath;

  get shouldShowPath() {
    const { showPath, scope, role } = this.args;
    return showPath || scope || role;
  }
}
