/**
 * Copyright (c) HashiCorp, Inc.
 * SPDX-License-Identifier: BUSL-1.1
 */

import { action } from '@ember/object';
import Component from '@glimmer/component';
import { tracked } from '@glimmer/tracking';

const pgpKeyFileDefault = () => ({ value: '' });

/**
 * @module ChoosePgpKeyForm
 * ChoosePgpKeyForm component is used for DR Operation Token Generation workflow. It provides
 * an interface for the user to upload or paste a PGP key for use
 *
 * @example
 * ```js
 * <ChoosePgpKeyForm @onCancel={{this.reset}} @onSubmit={{handleGenerateWithPgpKey}}>
 * ```
 * @param {function} onCancel - required - This function will be triggered when the modal intends to be closed
 * @param {function} onSubmit - required - When the PGP key is confirmed, it will call this method with the pgpKey value as the only param
 * @param {string} buttonText - Button text for onSubmit. Defaults to "Continue with key"
 * @param {string} formText - Form text above where the users uploads or pastes the key. Has default
 */
export default class ChoosePgpKeyForm extends Component {
  @tracked pgpKeyFile = pgpKeyFileDefault();
  @tracked selectedPgp = '';

  get pgpKey() {
    return this.pgpKeyFile.value;
  }

  get buttonText() {
    return this.args.buttonText || 'Continue with key';
  }

  get formText() {
    return (
      this.args.formText ||
      'Choose a PGP Key from your computer or paste the contents of one in the form below.'
    );
  }

  @action setKey(_, keyFile) {
    this.pgpKeyFile = keyFile;
  }

  // Form submit actions:
  @action usePgpKey(evt) {
    evt.preventDefault();
    this.selectedPgp = this.pgpKey;
  }
  @action handleSubmit(evt) {
    evt.preventDefault();
    this.args.onSubmit(this.pgpKey);
  }
}
