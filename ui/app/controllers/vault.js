/**
 * Copyright IBM Corp. 2016, 2025
 * SPDX-License-Identifier: BUSL-1.1
 */

import Controller from '@ember/controller';

export default class VaultController extends Controller {
  queryParams = [{ redirectTo: 'redirect_to' }];
  redirectTo = '';
}
