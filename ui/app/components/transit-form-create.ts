/**
 * Copyright IBM Corp. 2016, 2026
 * SPDX-License-Identifier: BUSL-1.1
 */

import Component from '@glimmer/component';
import { service } from '@ember/service';
import { action } from '@ember/object';
import { tracked } from '@glimmer/tracking';
import { HTMLElementEvent } from 'vault/forms';
import FlashMessageService from 'vault/services/flash-messages';
import ApiService from 'vault/services/api';
import RouterService from '@ember/routing/router-service';
import TransitKeyForm from 'vault/forms/transit/key';
import { ValidationMap } from 'vault/vault/app-types';
import { task } from 'ember-concurrency';

interface Args {
  form: TransitKeyForm;
}
export default class TransitFormCreate extends Component<Args> {
  @service declare readonly router: RouterService;
  @service declare readonly flashMessages: FlashMessageService;
  @service declare readonly api: ApiService;

  @tracked errorBanner = '';
  @tracked modelValidations: ValidationMap | null = null;

  @action
  handleAutoRotateChange(ttlObj: { enabled: boolean; goSafeTimeString: string }) {
    const { data } = this.args.form;
    if (ttlObj.enabled) {
      data.auto_rotate_period = ttlObj.goSafeTimeString;
    } else {
      data.auto_rotate_period = '0s';
    }
  }

  @task
  *save(event: HTMLElementEvent<HTMLFormElement>) {
    event.preventDefault();

    try {
      const { data, isValid, state } = this.args.form.toJSON();
      this.modelValidations = isValid ? null : state;

      if (isValid) {
        yield this.api.secrets.transitCreateKey(data['name'] as string, data['backend'] as string, data);

        this.flashMessages.success('Key successfully created.');
        this.router.transitionTo('vault.cluster.secrets.backend.show', data['backend'], data['name'], {
          queryParams: { tab: 'details' },
        });
      }
    } catch (error) {
      const { message } = yield this.api.parseError(error);
      this.errorBanner = message;
    }
  }
}
