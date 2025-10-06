/**
 * Copyright (c) HashiCorp, Inc.
 * SPDX-License-Identifier: BUSL-1.1
 */

import Controller, { inject as controller } from '@ember/controller';

import type ClientsCountsController from '../counts';

export default class ClientsCountsOverviewController extends Controller {
  @controller('vault.cluster.clients.counts') declare readonly countsController: ClientsCountsController;
}
