/**
 * Copyright (c) HashiCorp, Inc.
 * SPDX-License-Identifier: BUSL-1.1
 */

import ActivityComponent from '../activity';

export default class ClientsAcmePageComponent extends ActivityComponent {
  title = 'ACME usage';
  get description() {
    return 'ACME clients which interacted with Vault for the first time each month. Each bar represents the total new ACME clients for that month.';
  }
}
