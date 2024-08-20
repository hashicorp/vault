/**
 * Copyright (c) HashiCorp, Inc.
 * SPDX-License-Identifier: BUSL-1.1
 */

import ActivityComponent from '../activity';

export default class ClientsAcmePageComponent extends ActivityComponent {
  title = 'ACME usage';
  get description() {
    return `This data can be used to understand how many ACME clients have been used for the queried ${
      this.isDateRange ? 'date range' : 'month'
    }. Each ACME request is counted as one client.`;
  }
}
