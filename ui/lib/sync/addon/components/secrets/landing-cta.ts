/**
 * Copyright IBM Corp. 2016, 2025
 * SPDX-License-Identifier: BUSL-1.1
 */

import Component from '@glimmer/component';
import { service } from '@ember/service';

import type FlagsService from 'vault/services/flags';

interface Args {
  isActivated: boolean;
}

export default class LandingCtaComponent extends Component<Args> {
  @service declare readonly flags: FlagsService;
}
