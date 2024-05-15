/**
 * Copyright (c) HashiCorp, Inc.
 * SPDX-License-Identifier: BUSL-1.1
 */

import Component from '@glimmer/component';
import { service } from '@ember/service';
import { allMethods, methods } from 'vault/helpers/mountable-auth-methods';
import { allEngines, mountableEngines } from 'vault/helpers/mountable-secret-engines';

/**
 *
 * @module MountBackendTypeForm
 * MountBackendTypeForm components are used to display type options for
 * mounting either an auth method or secret engine.
 *
 * @example
 * ```js
 * <MountBackend::TypeForm @setMountType={{this.setMountType}} @mountType="secret" />
 * ```
 * @param {CallableFunction} setMountType - function will receive the mount type string. Should update the model type value
 * @param {string} [mountType=auth] - mount type can be `auth` or `secret`
 */

export default class MountBackendTypeForm extends Component {
  @service version;

  get secretEngines() {
    return this.version.isEnterprise ? allEngines() : mountableEngines();
  }

  get authMethods() {
    return this.version.isEnterprise ? allMethods() : methods();
  }

  get mountTypes() {
    return this.args.mountType === 'secret' ? this.secretEngines : this.authMethods;
  }
}
