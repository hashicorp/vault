/**
 * Copyright IBM Corp. 2016, 2026
 * SPDX-License-Identifier: BUSL-1.1
 */

import Component from '@glimmer/component';
import { action } from '@ember/object';
import { tracked } from '@glimmer/tracking';
import { service } from '@ember/service';
import { task } from 'ember-concurrency';

/**
 * MfaMethodForm component
 *
 * @example
 * ```js
 * <Mfa::MethodForm @model={{this.model}} @hasActions={{true}} @onSave={{this.onSave}} @onClose={{this.onClose}} />
 * ```
 * @param {Object} model - MFA method model
 * @param {boolean} [hasActions] - whether the action buttons will be rendered or not
 * @param {onSave} [onSave] - callback when save is successful
 * @param {onClose} [onClose] - callback when cancel is triggered
 */
export default class MfaMethodForm extends Component {
  @service flashMessages;
  @service api;

  @tracked editValidations;
  @tracked isEditModalActive = false;

  @task
  *save() {
    const { data } = this.args.form.toJSON();

    try {
      if (this.args.form.type === 'totp') {
        yield this.api.identity.mfaUpdateTotpMethod(data.id, { ...data });
      } else if (this.args.form.type === 'duo') {
        yield this.api.identity.mfaUpdateDuoMethod(data.id, { ...data });
      } else if (this.args.form.type === 'okta') {
        yield this.api.identity.mfaUpdateOktaMethod(data.id, { ...data });
      } else if (this.args.form.type === 'pingid') {
        yield this.api.identity.mfaUpdatePingIdaMethod(data.id, { ...data });
      }
      this.args.onSave();
    } catch (e) {
      this.flashMessages.danger(e.errors?.join('. ') || e.message);
    }
  }

  @action
  cancel() {
    this.args?.onClose?.();
  }

  @action
  async initSave(e) {
    e.preventDefault();
    const { isValid, state } = await this.args.form.toJSON();
    if (isValid) {
      this.isEditModalActive = true;
    } else {
      this.editValidations = state;
    }
  }
}
