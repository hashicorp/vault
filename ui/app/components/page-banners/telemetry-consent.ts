/**
 * Copyright IBM Corp. 2016, 2026
 * SPDX-License-Identifier: BUSL-1.1
 */

import Component from '@glimmer/component';
import { action } from '@ember/object';
import { service } from '@ember/service';

import type AnalyticsService from 'vault/services/analytics';

export default class TelemetryConsent extends Component {
  @service declare readonly analytics: AnalyticsService;

  @action
  recordConsent(userResponse: 'accept' | 'decline') {
    this.analytics.recordConsent(userResponse === 'accept');
  }
}
