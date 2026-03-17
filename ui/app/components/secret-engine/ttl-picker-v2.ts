/**
 * Copyright IBM Corp. 2016, 2025
 * SPDX-License-Identifier: BUSL-1.1
 */

import Component from '@glimmer/component';
import { action } from '@ember/object';
import { service } from '@ember/service';
import { tracked } from '@glimmer/tracking';
import { convertFromSeconds, durationToSeconds, largestUnitFromSeconds } from 'core/utils/duration-utils';

import type FlashMessageService from 'vault/services/flash-messages';
import type SecretsEngineResource from 'vault/resources/secrets/engine';
import type ApiService from 'vault/services/api';
import type RouterService from '@ember/routing/router-service';
import type { HTMLElementEvent } from 'vault/forms';

/**
 * @module TtlPickerV2 handles the display of the ttl picker fo the lease duration card in general settings.
 * 
 * @example
 * <SecretEngine::TtlPickerV2
    @model={{this.model}}
    @isDefaultTtlPicker={{boolean}}
    />
 *
 * @param {object} model - A model contains a secret engine resource, lease config from the sys/internal endpoint.
 * @param {boolean} isDefaultTtlPicker - isDefaultTtlPicker is a boolean that determines if the picker is default or max ttl.
 */

interface Args {
  model: {
    secretsEngine: SecretsEngineResource;
  };
  initialUnit: string;
  ttlKey: 'default_lease_ttl' | 'max_lease_ttl';
}

export default class TtlPickerV2 extends Component<Args> {
  systemDefaultTtl = 0;

  @service declare readonly flashMessages: FlashMessageService;
  @service declare readonly api: ApiService;
  @service declare readonly router: RouterService;

  @tracked selectedUnit = 's';
  @tracked time = '';
  @tracked errorMessage = '';

  constructor(owner: unknown, args: Args) {
    super(owner, args);

    this.initializeTtl();
  }

  initializeTtl() {
    const ttlValue = this.args?.model?.secretsEngine?.config[this.args.ttlKey];

    let seconds = 0;

    if (typeof ttlValue === 'number') {
      // if the passed value is a number, assume unit is seconds
      seconds = ttlValue;
    } else {
      const parseDuration = durationToSeconds(ttlValue || '');
      // if parsing fails leave it empty
      if (parseDuration === null) {
        this.time = ttlValue || '';
        this.selectedUnit = this.args.initialUnit;
        return;
      }
      seconds = parseDuration;
    }

    const unit = largestUnitFromSeconds(seconds);
    const time = convertFromSeconds(seconds, unit);
    this.time = time.toString() || '';
    this.selectedUnit = unit;
  }

  // reinitializes TTL when the model changes
  // after tuning, the model data is refreshed and the TTL units need to be recalculated & re-set
  @action
  onModelChange() {
    this.initializeTtl();
  }

  get unitOptions() {
    return [
      { label: 'seconds', value: 's' },
      { label: 'minutes', value: 'm' },
      { label: 'hours', value: 'h' },
      { label: 'days', value: 'd' },
    ];
  }

  get formField() {
    return {
      label:
        this.args?.ttlKey === 'default_lease_ttl'
          ? 'Default time-to-live (TTL)'
          : 'Maximum time-to-live (TTL)',
      helperText:
        this.args?.ttlKey === 'default_lease_ttl'
          ? 'How long secrets in this engine stay valid.'
          : 'Maximum extension for the secrets life beyond default.',
    };
  }

  @action
  setTtlTime(event: HTMLElementEvent<HTMLInputElement>) {
    this.errorMessage = '';
    if (isNaN(Number(event.target.value))) {
      this.errorMessage = 'Only use numbers for this setting.';
      return;
    }
    this.time = event.target.value;
    this.args.model.secretsEngine.config[this.args.ttlKey] = event.target.value;
  }

  @action
  setUnit(event: HTMLElementEvent<HTMLSelectElement>) {
    this.selectedUnit = event.target.value;
  }
}
