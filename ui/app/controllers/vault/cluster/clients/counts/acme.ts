/**
 * Copyright (c) HashiCorp, Inc.
 * SPDX-License-Identifier: BUSL-1.1
 */

import Controller, { inject as controller } from '@ember/controller';
import type ClientsCountsController from '../counts';

export default class ClientsCountsAcmeController extends Controller {
  // not sure why this needs to be cast to never but this definitely accepts a string to point to the controller
  @controller('vault.cluster.clients.counts' as never)
  declare readonly countsController: ClientsCountsController;
}
