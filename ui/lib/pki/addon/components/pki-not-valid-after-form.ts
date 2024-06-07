/**
 * Copyright (c) HashiCorp, Inc.
 * SPDX-License-Identifier: BUSL-1.1
 */

import Component from '@glimmer/component';
import { action } from '@ember/object';
import { tracked } from '@glimmer/tracking';
import { format } from 'date-fns';
import type { HTMLElementEvent } from 'forms';

/**
 * <PkiNotValidAfterForm /> components are used to manage two mutually exclusive role options in the form.
 */
interface Args {
  model: {
    notAfter: string;
    ttl: string | number;
    set: (key: string, value: string | number) => void;
  };
}

export default class PkiNotValidAfterForm extends Component<Args> {
  @tracked groupValue = 'ttl';
  @tracked cachedNotAfter: string;
  @tracked cachedTtl: string | number;
  @tracked formDate: string;

  constructor(owner: unknown, args: Args) {
    super(owner, args);
    const { model } = this.args;
    this.cachedNotAfter = model.notAfter || '';
    this.formDate = this.calculateFormDate(model.notAfter);
    this.cachedTtl = model.ttl || '';
    if (model.notAfter) {
      this.groupValue = 'specificDate';
    }
  }

  calculateFormDate(value: string) {
    // API expects and returns full ISO string
    // but the form input only accepts yyyy-MM-dd format
    if (value) {
      return format(new Date(value), 'yyyy-MM-dd');
    }
    return '';
  }

  @action onRadioButtonChange(selection: string) {
    this.groupValue = selection;
    // Clear the previous selection if they have clicked the other radio button.
    if (selection === 'specificDate') {
      this.args.model.ttl = '';
      this.args.model.notAfter = this.cachedNotAfter;
      this.formDate = this.calculateFormDate(this.cachedNotAfter);
    }
    if (selection === 'ttl') {
      this.args.model.notAfter = '';
      this.args.model.ttl = this.cachedTtl;
      this.formDate = '';
    }
  }

  @action setAndBroadcastTtl(ttlObject: { enabled: boolean; goSafeTimeString: string }) {
    const { enabled, goSafeTimeString } = ttlObject;
    if (this.groupValue === 'specificDate') {
      // do not save ttl on the model unless the ttl radio button is selected
      return;
    }
    const ttlVal = enabled === true ? goSafeTimeString : 0;
    this.cachedTtl = ttlVal;
    this.args.model.ttl = ttlVal;
  }

  @action setAndBroadcastInput(evt: HTMLElementEvent<HTMLInputElement>) {
    const setDate = evt.target.valueAsDate?.toISOString();
    if (!setDate) return;

    this.cachedNotAfter = setDate;
    this.args.model.notAfter = setDate;
    this.formDate = this.calculateFormDate(setDate);
  }
}
