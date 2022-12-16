import Component from '@glimmer/component';
import { inject as service } from '@ember/service';
import { methods } from 'vault/helpers/mountable-auth-methods';
import { allEngines, mountableEngines } from 'vault/helpers/mountable-secret-engines';
import { tracked } from '@glimmer/tracking';

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
 * @param {CallableFunction} setMountType - function will recieve the mount type string. Should update the model type value
 * @param {string} [mountType=auth] - mount type can be `auth` or `secret`
 */

export default class MountBackendTypeForm extends Component {
  @service version;
  @tracked selection;

  get secretEngines() {
    return this.version.isEnterprise ? allEngines() : mountableEngines();
  }

  get mountTypes() {
    return this.args.mountType === 'secret' ? this.secretEngines : methods();
  }
}
