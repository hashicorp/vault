/**
 * Copyright (c) HashiCorp, Inc.
 * SPDX-License-Identifier: BUSL-1.1
 */

import Component from '@ember/component';
import { service } from '@ember/service';
import UsageService from 'vault/services/usage';

export default class ClientsActivityComponent extends Component {
  @service declare readonly usage: UsageService;
}
