/**
 * Copyright IBM Corp. 2016, 2025
 * SPDX-License-Identifier: BUSL-1.1
 */

import Component from '@glimmer/component';
import { tracked } from '@glimmer/tracking';
import { action } from '@ember/object';
import { service } from '@ember/service';
import { normalizeMetricData, NormalizedBillingMetrics } from 'vault/utils/metrics-helpers';

import type ApiService from 'vault/services/api';
import type { Month, NormalizedMetricsData } from 'vault/vault/billing/overview';
import type { SystemReadBillingOverviewResponse } from '@hashicorp/vault-client-typescript';

const REFRESH_PERIOD_MS = 10 * 60 * 1000 + 30 * 1000; // 10 minutes 30 seconds

export default class BillingPageOverview extends Component {
  @service declare readonly api: ApiService;

  @tracked selectedDateOption: Month | null | undefined = null;
  @tracked normalizedMetricData: NormalizedMetricsData | undefined = {};
  @tracked months: Month[] = [];

  /** Reference to the scheduled timer, used to cancel on cleanup. */
  private _timer: ReturnType<typeof setTimeout> | null = null;

  /** Milliseconds to wait between each poll. Updated dynamically based on API response. */
  private _interval = 5000;

  detailsByMetric = {
    Secrets: [
      NormalizedBillingMetrics.STATIC_SECRETS_KV,
      NormalizedBillingMetrics.DYNAMIC_ROLES,
      NormalizedBillingMetrics.STATIC_ROLES,
    ],
    'Credential units': [
      NormalizedBillingMetrics.PKI_UNITS_TOTAL,
      NormalizedBillingMetrics.SSH_UNITS_OTP_UNITS,
      NormalizedBillingMetrics.SSH_UNITS_CERTIFICATE_UNITS,
    ],
    'Data protection calls': [
      NormalizedBillingMetrics.DATA_PROTECTION_CALLS_TRANSFORM,
      NormalizedBillingMetrics.DATA_PROTECTION_CALLS_TRANSIT,
    ],
    'Managed keys': [NormalizedBillingMetrics.MANAGED_KEYS_TOTP, NormalizedBillingMetrics.MANAGED_KEYS_KMSE],
  };

  constructor(owner: unknown, args: object) {
    super(owner, args);
    this.startPoll();
  }

  get selectedDate() {
    return this.selectedDateOption ?? this.months[0] ?? null;
  }

  /**
   * Calculates how long to wait before the next poll based on when the data was last updated.
   * Waits until 10m30s after `updated_at`, so polls align with the server's refresh cadence.
   */
  calculatePollingInterval(updatedAt: string): number {
    const msUntilRefresh = new Date(updatedAt).getTime() + REFRESH_PERIOD_MS - Date.now();
    // If data is already stale, wait a full period rather than re-polling the api immediately.
    return msUntilRefresh > 0 ? msUntilRefresh : REFRESH_PERIOD_MS;
  }

  fetchBillingMetrics = async () => {
    const response: SystemReadBillingOverviewResponse | null | undefined =
      await this.api.sys.systemReadBillingOverview();
    this.months = (response?.months as Month[]) || [];
    const updatedMonthFromSelectedMonth = this.months.find(
      (month: Month) => month.month === this.selectedDateOption?.month
    );
    const updatedMonth: Month | undefined = updatedMonthFromSelectedMonth || this.months[0];

    if (updatedMonth?.updated_at) {
      this._interval = this.calculatePollingInterval(updatedMonth.updated_at);
    }

    this.onDateChange(updatedMonth ?? null);
    return this.months;
  };

  /**
   * Starts the polling loop, invoking fetchBillingMetrics immediately and then
   * repeatedly on each interval. No-ops if polling is already active.
   */
  startPoll() {
    if (this._timer) return;

    const poll = async () => {
      try {
        await this.fetchBillingMetrics();
      } catch (e) {
        // Error fetching billing metrics
      } finally {
        // Schedule the next poll using the current interval value,
        // which may have been updated by the callback.
        this._timer = setTimeout(poll, this._interval);
      }
    };

    poll();
  }

  /**
   * Stops the polling loop and cancels any pending scheduled poll.
   */
  stopPoll() {
    if (this._timer) {
      clearTimeout(this._timer);
      this._timer = null;
    }
  }

  metricsForCard = (cardData: string[]) => {
    const metrics: NormalizedMetricsData = {};
    // Iterate over keys for that card's data
    // so only relevant metrics are passed to each card
    for (const key of cardData) {
      metrics[key] = this.normalizedMetricData?.[key];
    }

    return metrics;
  };

  @action
  onDateChange(dropdownOption: Month | null | undefined) {
    this.selectedDateOption = dropdownOption;
    this.normalizedMetricData = normalizeMetricData(dropdownOption);
  }

  willDestroy() {
    super.willDestroy();
    this.stopPoll();
  }
}
