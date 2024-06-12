/**
 * Copyright (c) HashiCorp, Inc.
 * SPDX-License-Identifier: BUSL-1.1
 */

import { service } from '@ember/service';
import Controller from '@ember/controller';
import config from '../config/environment';

export default class ApplicationController extends Controller {
  @service auth;
  @service store;
  env = config.environment;
}
