/**
 * Copyright IBM Corp. 2016, 2025
 * SPDX-License-Identifier: BUSL-1.1
 */

import Component from '@glimmer/component';
import { action } from '@ember/object';
import { service } from '@ember/service';
import { tracked } from '@glimmer/tracking';
import { task } from 'ember-concurrency';

import type { InternalClientActivityReadConfigurationResponse } from '@hashicorp/vault-client-typescript';
import type RouterService from '@ember/routing/router-service';
import type { HTMLElementEvent } from 'vault/forms';
import type ApiService from 'vault/services/api';
import type Owner from '@ember/owner';

interface Args {
  config: InternalClientActivityReadConfigurationResponse;
  mode: 'show' | 'edit';
}

export default class ConfigComponent extends Component<Args> {
  @service declare readonly router: RouterService;
  @service declare readonly api: ApiService;

  @tracked modalOpen = false;
  @tracked declare enabled: boolean;
  @tracked validationError = '';
  @tracked errorMessage = '';

  constructor(owner: Owner, args: Args) {
    super(owner, args);
    const { enabled = '' } = args.config;
    // possible config values are 'enable', 'disable', 'default-enabled', 'default-disabled'
    this.enabled = enabled.includes('enable');
  }

  get infoRows() {
    return [
      {
        label: 'Usage data collection',
        helperText: 'Enable or disable collecting data to track clients.',
        value: this.enabled ? 'On' : 'Off',
      },
      {
        label: 'Retention period',
        helperText: 'The number of months of activity logs to maintain for client tracking.',
        value: this.args.config.retention_months,
      },
    ];
  }

  get modalTitle() {
    return `Turn usage tracking ${this.enabled ? 'on' : 'off'}?`;
  }

  @action
  onSubmit(event: HTMLElementEvent<HTMLFormElement>) {
    event.preventDefault();

    this.validationError = '';
    // since minimum_retention_months may be returned as 0, default to 48 which is the documented minimum
    // https://developer.hashicorp.com/vault/api-docs/system/internal-counters#retention_months
    const { minimum_retention_months, retention_months, enabled = '' } = this.args.config;
    const minRetention = minimum_retention_months || 48;

    if (Number(retention_months) < minRetention) {
      this.validationError = `Retention period must be greater than or equal to ${minRetention}.`;
    } else if (Number(retention_months) > 60) {
      this.validationError = 'Retention period must be less than or equal to 60.';
    }
    // if form is valid and enabled value has changed show the confirmation modal
    // values for enabled may include 'default-' so check for inclusion of enable or disable
    if (!this.validationError) {
      const didChange = enabled.includes('enable') ? !this.enabled : !!this.enabled;
      if (didChange) {
        // the modal confirm action will trigger the save task directly
        this.modalOpen = true;
      } else {
        this.save.perform();
      }
    }
  }

  save = task(async () => {
    try {
      const payload = {
        enabled: this.enabled ? 'enable' : 'disable',
        retention_months: Number(this.args.config.retention_months),
      };
      await this.api.sys.internalClientActivityConfigure(payload);
      this.router.transitionTo('vault.cluster.clients.config');
    } catch (error) {
      const { message } = await this.api.parseError(error);
      this.errorMessage = message;
    }
  });
}
