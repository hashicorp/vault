/**
 * Copyright (c) HashiCorp, Inc.
 * SPDX-License-Identifier: BUSL-1.1
 */

import { action } from '@ember/object';
import Component from '@glimmer/component';
import { tracked } from '@glimmer/tracking';
import { HTMLElementEvent } from 'vault/forms';
import SecretsEngineResource from 'vault/resources/secrets/engine';
interface Args {
  model: SecretsEngineResource;
}

export default class LeaseDuration extends Component<Args> {
  // TODO: When wiring up to parent, address variable names and usage, update onchange functions etc - reference ttl-picker.js
  @tracked enableTTL = false;
  @tracked enableMaxTTL = false;
  @tracked time = '';
  @tracked maxTime = '';
  @tracked unit = 's';
  @tracked maxUnit = 's';

  constructor(owner: unknown, args: Args) {
    super(owner, args);
  }

  @action
  setTTLType(event: HTMLElementEvent<HTMLInputElement>) {
    if (event.target.value === 'Custom') {
      this.enableTTL = true;
    } else {
      this.enableTTL = false;
    }
  }

  @action
  setMaxTTLType(event: HTMLElementEvent<HTMLInputElement>) {
    if (event.target.value === 'Custom') {
      this.enableMaxTTL = true;
    } else {
      this.enableMaxTTL = false;
    }
  }

  @action
  setTtlTime(event: HTMLElementEvent<HTMLInputElement>) {
    this.time = event.target.value;
  }

  @action
  setMaxTtlTime(event: HTMLElementEvent<HTMLInputElement>) {
    this.maxTime = event.target.value;
  }

  @action
  setUnit(event: HTMLElementEvent<HTMLSelectElement>) {
    this.unit = event.target.value;
  }

  @action
  setMaxUnit(event: HTMLElementEvent<HTMLSelectElement>) {
    this.maxUnit = event.target.value;
  }
}
