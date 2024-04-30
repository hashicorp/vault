/**
 * Copyright (c) HashiCorp, Inc.
 * SPDX-License-Identifier: BUSL-1.1
 */

import Component from '@glimmer/component';

interface Args {
  showSecretsSync: boolean;
}

export default class RunningTotal extends Component<Args> {
  get chartContainerText() {
    const { showSecretsSync } = this.args;
    return `The total clients in the specified date range, displayed per month. This includes entity, non-entity${
      showSecretsSync ? ', ACME and secrets sync clients' : ' and ACME clients'
    }. The total client count number is an important consideration for Vault billing.`;
  }
}
