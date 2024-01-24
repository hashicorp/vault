/**
 * Copyright (c) HashiCorp, Inc.
 * SPDX-License-Identifier: BUSL-1.1
 */

import ActivityComponent from '../activity';

export default class ClientsTokenPageComponent extends ActivityComponent {
  legend = [
    { key: 'entity_clients', label: 'entity clients' },
    { key: 'non_entity_clients', label: 'non-entity clients' },
  ];

  get hasAverageNewClients() {
    return (
      typeof this.average(this.byMonthNewClients, 'entity_clients') === 'number' ||
      typeof this.average(this.byMonthNewClients, 'non_entity_clients') === 'number'
    );
  }
}
