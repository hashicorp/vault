/**
 * Copyright (c) HashiCorp, Inc.
 * SPDX-License-Identifier: BUSL-1.1
 */

import Component from '@glimmer/component';

interface Args {
  isSecretsSyncActivated: boolean;
}

export default class RunningTotal extends Component<Args> {
  get chartContainerText() {
    const { isSecretsSyncActivated } = this.args;
    return `The total clients in the specified date range, displayed per month. This includes entity, non-entity${
      isSecretsSyncActivated ? ', ACME and secrets sync clients' : ' and ACME clients'
    }. The total client count number is an important consideration for Vault billing.`;
  }
}
