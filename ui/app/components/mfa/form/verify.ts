/**
 * Copyright (c) HashiCorp, Inc.
 * SPDX-License-Identifier: BUSL-1.1
 */

import Component from '@glimmer/component';
import { action } from '@ember/object';
import { numberToWord } from 'vault/helpers/number-to-word';

import type { MfaConstraintState } from 'vault/vault/auth/mfa';
import type { HTMLElementEvent } from 'vault/forms';

interface Args {
  codeDelayMessage: string;
  constraints: MfaConstraintState[];
  countdown: string;
  error: string;
  isLoading: boolean;
  methodAlreadyEnrolled: CallableFunction;
  onCancel: CallableFunction;
  onSelect: CallableFunction;
  onVerify: CallableFunction;
}

export default class MfaFormVerify extends Component<Args> {
  get description() {
    const base = 'Multi-factor authentication is enabled for your account.';
    if (this.args.constraints.length > 1) {
      const num = this.args.constraints.length;
      const word = num === 1 ? 'method' : 'methods';
      return base + ` ${numberToWord(num, true)} ${word} are required for successful authentication.`;
    }
    if (this.singleConstraint?.selectedMethod?.uses_passcode) {
      return base + ' Enter your authentication code to log in.';
    }
    return base;
  }

  get singleConstraint() {
    return this.args.constraints.length === 1 ? this.args.constraints[0] : null;
  }

  @action
  handleSubmit(e: HTMLElementEvent<HTMLFormElement>) {
    e.preventDefault();
    this.args.onVerify();
  }

  // Template helper
  sortConstraints = (constraints: MfaConstraintState[]) => {
    const userInteraction = constraints.filter((c) => !c.selectedMethod);
    const others = constraints.filter((c) => c.selectedMethod);
    return [...userInteraction, ...others];
  };

  // Even if multiple methods fail, the API seems to only ever return whichever failed first.
  // Sample error message:
  // 'login MFA validation failed for methodID: [9e953c14-9d8e-7443-079e-4b13723e2aef]'
  hasValidationError = (id: string) => this.args.error?.includes(id);
}
