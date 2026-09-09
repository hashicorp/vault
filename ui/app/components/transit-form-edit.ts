/**
 * Copyright IBM Corp. 2016, 2025
 * SPDX-License-Identifier: BUSL-1.1
 */

import Component from '@glimmer/component';
import { service } from '@ember/service';
import { action } from '@ember/object';
import { tracked } from '@glimmer/tracking';
import ApiService from 'vault/services/api';
import FlashMessageService from 'vault/services/flash-messages';
import RouterService from '@ember/routing/router-service';
import TransitKeyForm from 'vault/forms/transit/key';
import { ValidationMap } from 'vault/vault/app-types';
import { HTMLElementEvent } from 'vault/forms';

type TransitKeyConfig = {
  name?: string;
  backend?: string;
  auto_rotate_period?: string;
  min_decryption_version?: number;
  min_encryption_version?: number;
  deletion_allowed?: boolean;
};

interface Args {
  form: TransitKeyForm;
}
export default class TransitFormEdit extends Component<Args> {
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

  @action async deleteKey() {
    const { backend, id } = this.args.form.data;
    try {
      await this.api.secrets.transitDeleteKey(id as string, backend as string);
      this.flashMessages.success(`'${id}' key deleted.`);
      this.router.transitionTo('vault.cluster.secrets.backend.list-root');
    } catch (error) {
      const { message } = await this.api.parseError(error);
      this.errorBanner = message;
    }
  }

  @action async save(event: HTMLElementEvent<HTMLFormElement>) {
    event.preventDefault();

    try {
      const { data, isValid, state } = this.args.form.toJSON();
      this.modelValidations = isValid ? null : state;

      if (isValid) {
        const {
          name,
          backend,
          auto_rotate_period,
          min_decryption_version,
          min_encryption_version,
          deletion_allowed,
        }: TransitKeyConfig = data;

        await this.api.secrets.transitConfigureKey(name as string, backend as string, {
          auto_rotate_period,
          min_decryption_version,
          min_encryption_version,
          deletion_allowed,
        });

        this.flashMessages.success('Key successfully updated.');
        this.router.transitionTo('vault.cluster.secrets.backend.show', backend, name, {
          queryParams: { tab: 'details' },
        });
      }
    } catch (error) {
      const { message } = await this.api.parseError(error);
      this.errorBanner = message;
    }
  }
}
