/**
 * Copyright (c) HashiCorp, Inc.
 * SPDX-License-Identifier: BUSL-1.1
 */

import Component from '@glimmer/component';

interface Args {
  isSecretsSyncActivated: boolean;
}

export default class VerticalBarChart extends Component<Args> {
  get chartContainerText() {
    const { isSecretsSyncActivated } = this.args;
    const prefix = 'The total clients in the specified date range. This includes entity';

    const mid = isSecretsSyncActivated ? ', non-entity, and secrets sync clients' : ' and non-entity clients';
    const suffix = '. The total client count number is an important consideration for Vault billing.';

    return `${prefix}${mid}${suffix}`;
  }
}
