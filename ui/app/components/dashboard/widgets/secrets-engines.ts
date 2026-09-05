/**
 * Copyright IBM Corp. 2016, 2026
 * SPDX-License-Identifier: BUSL-1.1
 */

import Component from '@glimmer/component';

import type SecretsEngineResource from 'vault/resources/secrets/engine';

import { service } from '@ember/service';
import { action } from '@ember/object';
import {
  DASHBOARD_SECRETS_ENGINE_VIEW_ALL,
  DASHBOARD_SECRETS_ENGINE_VIEW,
} from 'vault/utils/analytic-events';

import type AnalyticsService from 'vault/services/analytics';

/**
 * @module DashboardWidgetsSecretsEngines
 * DashboardWidgetsSecretsEngines component are used to display 5 secrets engines to the user.
 *
 * @example
 * ```js
 * <DashboardWidgetsSecretsEngines @secretsEngines={{@model.secretsEngines}} />
 * ```
 * @param {array} secretsEngines - list of secrets engines
 */

interface Args {
  secretsEngines: SecretsEngineResource[];
}

export default class DashboardWidgetsSecretsEngines extends Component<Args> {
  @service declare readonly analytics: AnalyticsService;

  get filteredSecretsEngines() {
    return this.args.secretsEngines?.filter((secretEngine) => secretEngine.shouldIncludeInList);
  }

  get firstFiveSecretsEngines() {
    return this.filteredSecretsEngines?.slice(0, 5);
  }

  @action
  async trackViewAll() {
    this.analytics.trackEvent(DASHBOARD_SECRETS_ENGINE_VIEW_ALL, {
      namespace: 'dashboard',
      action: 'clicked',
      elementId: 'dashboard-secrets-engine-view-all',
      channel: 'webpage',
    });
  }

  @action
  async trackViewEngine(backendType: string) {
    this.analytics.trackEvent(DASHBOARD_SECRETS_ENGINE_VIEW, {
      namespace: 'dashboard',
      action: 'clicked',
      elementId: 'dashboard-secrets-engine-view',
      channel: 'webpage',
      objectType: 'secrets-engine',
      object: backendType,
    });
  }
}
