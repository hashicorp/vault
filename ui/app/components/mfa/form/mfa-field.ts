/**
 * Copyright (c) HashiCorp, Inc.
 * SPDX-License-Identifier: BUSL-1.1
 */

import { action } from '@ember/object';
import Component from '@glimmer/component';

import type { HTMLElementEvent } from 'vault/forms';
import type MfaConstraint from 'vault/resources/mfa/constraint';

interface Args {
  constraint: MfaConstraint;
  onSelect: CallableFunction;
}

export default class MfaFormMfaField extends Component<Args> {
  @action
  setPasscode(constraint: MfaConstraint, e: HTMLElementEvent<HTMLInputElement>) {
    const { value } = e.target;
    constraint.setPasscode(value);
  }

  @action
  handleSelect(constraint: MfaConstraint, e: HTMLElementEvent<HTMLInputElement>) {
    const { value: id } = e.target;
    this.args.onSelect(constraint, id);
  }
}
