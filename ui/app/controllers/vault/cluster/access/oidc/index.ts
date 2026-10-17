/**
 * Copyright IBM Corp. 2016, 2025
 * SPDX-License-Identifier: BUSL-1.1
 */

import Controller from '@ember/controller';
import { service } from '@ember/service';

import type ThemeService from 'vault/services/theme';

export default class OidcIndexController extends Controller {
  @service declare readonly theme: ThemeService;
}
