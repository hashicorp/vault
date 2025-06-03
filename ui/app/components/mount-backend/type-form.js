/**
 * Copyright (c) HashiCorp, Inc.
 * SPDX-License-Identifier: BUSL-1.1
 */

import Component from '@glimmer/component';
import { service } from '@ember/service';
import { ALL_ENGINES } from 'vault/utils/all-engines-metadata';

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
    return this.version.isEnterprise
      ? ALL_ENGINES.map((engine) => engine.mountType !== 'auth')
      : ALL_ENGINES.map((engine) => engine.mountType !== 'auth' && !engine.requiresEnterprise);
  }

  get authMethods() {
    return this.version.isEnterprise
      ? ALL_ENGINES.map((engine) => engine.mountType !== 'secret')
      : ALL_ENGINES.map((engine) => engine.mountType !== 'secret' && !engine.requiresEnterprise);
  }

  get mountTypes() {
    return this.args.mountType === 'secret' ? this.secretEngines : this.authMethods;
  }
}
