/**
 * Copyright (c) HashiCorp, Inc.
 * SPDX-License-Identifier: BUSL-1.1
 */

import Component from '@glimmer/component';

import type { MfaConstraintState } from 'vault/vault/auth/mfa';

interface Args {
  constraints: MfaConstraintState[];
  onSelect: CallableFunction;
}

const METHOD_MAP = {
  totp: { label: 'TOTP', icon: 'history' },
  duo: { label: 'Duo', icon: 'duo-color' },
  okta: { label: 'Okta', icon: 'okta-color' },
  pingid: { label: 'PingID', icon: 'ping-identity-color' },
};

export default class MfaFormChooseMethod extends Component<Args> {
  get nonSelfEnrollMethods() {
    return this.singleConstraint?.methods.filter((m) => !m.self_enrollment_enabled);
  }

  get selfEnrollConstraints() {
    return this.args.constraints.filter((c) => !!c.selfEnrollMethod);
  }

  get singleConstraint() {
    // Prioritize self-enroll constraints so the user can setup TOTP before moving onto validating.
    if (this.selfEnrollConstraints?.length === 1) {
      return this.selfEnrollConstraints[0];
    }
    if (this.args.constraints.length === 1) {
      return this.args.constraints[0];
    }
    return null;
  }

  // TEMPLATE HELPERS
  displayIcon = (methodType: 'duo' | 'okta' | 'totp' | 'pingid') => METHOD_MAP[methodType].icon;
  displayLabel = (methodType: 'duo' | 'okta' | 'totp' | 'pingid') => METHOD_MAP[methodType].label;
}
