/**
 * Copyright IBM Corp. 2016, 2025
 * SPDX-License-Identifier: BUSL-1.1
 */

import Controller from '@ember/controller';

export default class PkiExternalOrdersIndexController extends Controller {
  queryParams = ['within'];
  within = null;
}
