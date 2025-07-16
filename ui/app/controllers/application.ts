/**
 * Copyright (c) HashiCorp, Inc.
 * SPDX-License-Identifier: BUSL-1.1
 */

import { service } from '@ember/service';
import Controller from '@ember/controller';
import config from '../config/environment';

import type AuthService from 'vault/vault/services/auth';
import type Store from '@ember-data/store';

export default class ApplicationController extends Controller {
  @service declare readonly auth: AuthService;
  @service declare readonly store: Store;
  env = config.environment;
}
