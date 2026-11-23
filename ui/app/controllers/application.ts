/**
 * Copyright IBM Corp. 2016, 2025
 * SPDX-License-Identifier: BUSL-1.1
 */

import { service } from '@ember/service';
import Controller from '@ember/controller';
import config from '../config/environment';

import type AuthService from 'vault/vault/services/auth';

export default class ApplicationController extends Controller {
  @service declare readonly auth: AuthService;
  env = config.environment;
}
