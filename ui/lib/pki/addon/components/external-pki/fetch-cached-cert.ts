/**
 * Copyright IBM Corp. 2016, 2026
 * SPDX-License-Identifier: BUSL-1.1
 */

import { action } from '@ember/object';
import { durationToSeconds } from 'core/utils/duration-utils';
import { number } from 'vault/utils/forms/validators';
import { task } from 'ember-concurrency';
import { tracked } from '@glimmer/tracking';
import { validate } from 'vault/utils/forms/validate';
import Component from '@glimmer/component';

import type { HTMLElementEvent } from 'vault/forms';
import type { HTTPQuery } from '@hashicorp/vault-client-typescript';
import type { ValidationMap, Validations } from 'vault/vault/app-types';

type ValidityType = 'duration' | 'percentage';

interface Args {
  onSubmit: CallableFunction;
  errorMessage: string;
}

export default class ExternalPkiFetchCachedCertComponent extends Component<Args> {
  @tracked durationUnit = 's';
  @tracked identifiers = '';
  @tracked minValidityValue = '';
  @tracked validationErrors: ValidationMap | null = null;
  @tracked validityType: ValidityType = 'duration';

  unitOptions = [
    { label: 'seconds', value: 's' },
    { label: 'minutes', value: 'm' },
    { label: 'hours', value: 'h' },
    { label: 'days', value: 'd' },
  ];

  get minValidityValidation() {
    const validator = ({ minValidityValue }: { minValidityValue: string }) => {
      if (!number(minValidityValue)) {
        return false;
      }

      const value = Number(minValidityValue);
      return this.validityType === 'duration' ? value > 0 : value >= 1 && value <= 100;
    };

    return [
      {
        validator,
        message:
          this.validityType === 'duration'
            ? 'Duration must be a positive number.'
            : 'Percentage must be a number between 1 and 100.',
      },
    ];
  }

  get validations(): Validations {
    return {
      identifiers: [{ type: 'presence', message: 'You must provide at least one identifier.' }],
      minValidityValue: [
        { type: 'presence', message: 'A minimum validity is required.' },
        ...this.minValidityValidation,
      ],
    };
  }

  @action
  handleInput(e: HTMLElementEvent<HTMLInputElement | HTMLSelectElement>) {
    const { name, value } = e.target;

    switch (name) {
      case 'durationUnit':
      case 'identifiers':
      case 'minValidityValue':
        this[name] = value;
        break;
      default:
        break;
    }
  }

  @action
  handleRadioSelect(event: HTMLElementEvent<HTMLInputElement>) {
    this.validityType = event.target.value as ValidityType;
    this.validationErrors = null;
  }

  @action
  onSubmit(event: HTMLElementEvent<HTMLFormElement>) {
    event.preventDefault();
    this.validationErrors = null;

    const formData = {
      identifiers: this.identifiers,
      minValidityValue: this.minValidityValue,
    };
    const { isValid, state } = validate(formData, this.validations);
    if (!isValid) {
      this.validationErrors = state;
      return;
    }

    // Map form data to API params
    const payload: HTTPQuery = { identifiers: this.identifiers };
    if (this.validityType === 'duration') {
      payload['min_validity_duration'] = durationToSeconds(`${this.minValidityValue}${this.durationUnit}`);
    } else {
      payload['min_validity_percentage'] = Number(this.minValidityValue.trim());
    }

    this.save.perform(payload);
  }

  save = task(async (payload) => {
    await this.args.onSubmit(payload);
  });

  // Template helpers

  validationError = (param: string): string => {
    const { isValid, errors } = this.validationErrors?.[param] ?? {};
    return !isValid && errors ? errors.join(' ') : '';
  };
}
