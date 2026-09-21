/**
 * Copyright IBM Corp. 2016, 2025
 * SPDX-License-Identifier: BUSL-1.1
 */
import Controller from '@ember/controller';
import { tracked } from '@glimmer/tracking';

export default class IdentityGroupsIndexController extends Controller {
  queryParams = ['page', 'pageSize'];
  @tracked page = 1;
  @tracked pageSize = 10;
}
