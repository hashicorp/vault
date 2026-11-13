/**
 * Copyright IBM Corp. 2026
 * SPDX-License-Identifier: BUSL-1.1
 */

import Component from '@glimmer/component';
import { service } from '@ember/service';
import { tracked } from '@glimmer/tracking';
import { task } from 'ember-concurrency';
import { waitFor } from '@ember/test-waiters';

import type ApiService from 'vault/services/api';
import type FlashMessageService from 'vault/services/flash-messages';
import type ControlGroupConfigForm from 'vault/forms/control-group-config';
import type { HTMLElementEvent } from 'vault/forms';

interface Args {
  form: ControlGroupConfigForm;
}

/**
 * @module ControlGroupConfigForm
 * Form for editing the cluster-wide Control Group configuration.
 *
 * @param {ControlGroupConfigForm} form - form instance holding the max_ttl value
 */
export default class ControlGroupConfigFormComponent extends Component<Args> {
  @service declare readonly api: ApiService;
  @service declare readonly flashMessages: FlashMessageService;

  @tracked errorBanner = '';

  save = task(
    waitFor(async (event: HTMLElementEvent<HTMLFormElement>) => {
      event.preventDefault();
      this.errorBanner = '';
      try {
        const { data } = this.args.form.toJSON();
        await this.api.sys.systemWriteConfigControlGroup(data);
        this.flashMessages.success('The configuration was saved successfully.');
      } catch (error) {
        const { message } = await this.api.parseError(error);
        this.errorBanner = message;
      }
    })
  );
}
