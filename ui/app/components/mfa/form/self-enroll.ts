/**
 * Copyright (c) HashiCorp, Inc.
 * SPDX-License-Identifier: BUSL-1.1
 */

import { action } from '@ember/object';
import Component from '@glimmer/component';
import { tracked } from '@glimmer/tracking';
import { HTMLElementEvent } from 'vault/forms';

import type { MfaConstraintState } from 'vault/vault/auth/mfa';

interface Args {
  constraints: MfaConstraintState[];
  onVerify: CallableFunction;
  onCancel: CallableFunction;
}

export default class MfaFormSelfEnroll extends Component<Args> {
  @tracked hasScannedQrCode = false;

  get description() {
    return this.hasScannedQrCode
      ? 'To verify your device, enter the code generated from your authenticator.'
      : 'Scan the QR code with your authenticator app. If you currently do not have a device on hand, you can copy the MFA secret below and enter it manually.';
  }

  get selfEnrollConstraint() {
    // Find the constraint with the QR code, only one will have one at a time
    return this.args.constraints.find((c) => !!c.qrCode);
  }

  @action
  handleSubmit(e: HTMLElementEvent<HTMLFormElement>) {
    e.preventDefault();
    // Clear out the QR Code
    const constraint = this.findConstraint();
    if (constraint) {
      constraint.qrCode = '';
    }

    this.args.onVerify();
  }

  private findConstraint = () =>
    this.args.constraints.find((c) => c.name === this.selfEnrollConstraint?.name);
}
