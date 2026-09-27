/**
 * Copyright IBM Corp. 2016, 2025
 * SPDX-License-Identifier: BUSL-1.1
 */

import Controller from '@ember/controller';
import { tracked } from '@glimmer/tracking';

export default class IdentityGroupAliasesIndexController extends Controller {
  // page refreshes the model (re-fetches from API) on change.
  // pageSize is client-side only — no model reload needed, but must survive route transitions.
  queryParams = ['page', 'pageSize'];
  @tracked page = 1;
  @tracked pageSize = 10;
}
