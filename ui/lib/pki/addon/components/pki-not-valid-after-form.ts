/**
 * Copyright IBM Corp. 2016, 2025
 * SPDX-License-Identifier: BUSL-1.1
 */

import Component from '@glimmer/component';
import { action } from '@ember/object';
import { tracked } from '@glimmer/tracking';
import { format } from 'date-fns';

import type { HTMLElementEvent } from 'forms';
import type PkiConfigGenerateForm from 'vault/forms/secrets/pki/config/generate';
import type PkiCertificateIssueForm from 'vault/forms/secrets/pki/certificate';
import type PkiIssuersSignIntermediateForm from 'vault/forms/secrets/pki/issuers/sign-intermediate';

/**
 * <PkiNotValidAfterForm /> components are used to manage two mutually exclusive role options in the form.
 */
interface Args {
  form: PkiConfigGenerateForm | PkiIssuersSignIntermediateForm | PkiCertificateIssueForm;
}

export default class PkiNotValidAfterForm extends Component<Args> {
  @tracked groupValue = 'ttl';
  @tracked cachedNotAfter: string;
  @tracked cachedTtl: string | number;
  @tracked formDate: string;

  constructor(owner: unknown, args: Args) {
    super(owner, args);
    const { not_after, ttl } = this.args.form.data;
    this.cachedNotAfter = not_after || '';
    this.formDate = this.calculateFormDate(this.cachedNotAfter);
    this.cachedTtl = ttl || '';
    if (this.cachedNotAfter) {
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
    const { data } = this.args.form;
    if (selection === 'specificDate') {
      data.ttl = '';
      data.not_after = this.cachedNotAfter;
      this.formDate = this.calculateFormDate(this.cachedNotAfter);
    }
    if (selection === 'ttl') {
      data.not_after = '';
      data.ttl = `${this.cachedTtl}`;
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
    this.args.form.data.ttl = ttlVal.toString();
  }

  @action setAndBroadcastInput(evt: HTMLElementEvent<HTMLInputElement>) {
    const setDate = evt.target.valueAsDate?.toISOString();
    if (!setDate) return;

    this.cachedNotAfter = setDate;
    this.args.form.data.not_after = setDate;
    this.formDate = this.calculateFormDate(setDate);
  }
}
