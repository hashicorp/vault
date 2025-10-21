/**
 * Copyright IBM Corp. 2016, 2025
 * SPDX-License-Identifier: BUSL-1.1
 */

import Controller, { inject as controller } from '@ember/controller';

import type ClientsCountsController from '../counts';

export default class ClientsCountsClientListController extends Controller {
  @controller('vault.cluster.clients.counts') declare readonly countsController: ClientsCountsController;
}
