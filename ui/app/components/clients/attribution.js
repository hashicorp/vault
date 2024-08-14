/**
 * Copyright (c) HashiCorp, Inc.
 * SPDX-License-Identifier: BUSL-1.1
 */

import Component from '@glimmer/component';
import { parseAPITimestamp } from 'core/utils/date-formatters';
import { isSameMonth } from 'date-fns';

/**
 * @module Attribution
 * Attribution components display the top 10 total client counts for namespaces or auth methods (mounts) during a billing period.
 * A horizontal bar chart shows on the right, with the top namespace/auth method and respective client totals on the left.
 *
 * @example
 *  <Clients::Attribution
 *    @totalClientAttribution={{this.totalClientAttribution}}
 *    @newClientAttribution={{this.newClientAttribution}}
 *    @selectedNamespace={{this.selectedNamespace}}
 *    @startTimestamp={{this.startTime}}
 *    @endTimestamp={{this.endTime}}
 *    @isHistoricalMonth={{false}}
 *    @responseTimestamp={{this.responseTimestamp}}
 *  />
 *
 * @param {array} totalClientAttribution - array of objects containing a label and breakdown of client counts for total clients
 * @param {array} newClientAttribution - array of objects containing a label and breakdown of client counts for new clients
 * @param {string} selectedNamespace - namespace selected from filter bar
 * @param {string} startTimestamp - timestamp string from activity response to render start date for CSV modal and whether copy reads 'month' or 'date range'
 * @param {string} endTimestamp - timestamp string from activity response to render end date for CSV modal and whether copy reads 'month' or 'date range'
 * @param {string} responseTimestamp -  ISO timestamp created in serializer to timestamp the response, renders in bottom left corner below attribution chart
 * @param {boolean} isHistoricalMonth - when true data is from a single, historical month so side-by-side charts should display for attribution data
 * @param {boolean} isSecretsSyncActivated - boolean to determine if secrets sync is activated
 */

export default class Attribution extends Component {
  get attributionLegend() {
    const attributionLegend = [
      { key: 'entity_clients', label: 'entity clients' },
      { key: 'non_entity_clients', label: 'non-entity clients' },
      { key: 'acme_clients', label: 'ACME clients' },
    ];

    if (this.args.isSecretsSyncActivated) {
      attributionLegend.push({ key: 'secret_syncs', label: 'secrets sync clients' });
    }
    return attributionLegend;
  }

  get isSingleMonth() {
    if (!this.args.startTimestamp && !this.args.endTimestamp) return false;
    const startDateObject = parseAPITimestamp(this.args.startTimestamp);
    const endDateObject = parseAPITimestamp(this.args.endTimestamp);
    return isSameMonth(startDateObject, endDateObject);
  }

  get isSingleNamespace() {
    // if a namespace is selected, then we're viewing top 10 auth methods (mounts)
    return !!this.args.selectedNamespace;
  }

  // truncate data before sending to chart component
  get barChartTotalClients() {
    return this.args.totalClientAttribution?.slice(0, 10);
  }

  get barChartNewClients() {
    return this.args.newClientAttribution?.slice(0, 10);
  }

  get topClientCounts() {
    // get top namespace or auth method
    return this.args.totalClientAttribution ? this.args.totalClientAttribution[0] : null;
  }

  get attributionBreakdown() {
    // display text for hbs
    return this.isSingleNamespace ? 'auth method' : 'namespace';
  }

  get chartText() {
    if (!this.args.totalClientAttribution) {
      return { description: 'There is a problem gathering data' };
    }
    const dateText = this.isSingleMonth ? 'month' : 'date range';
    switch (this.isSingleNamespace) {
      case true:
        return {
          description:
            'This data shows the top ten authentication methods by client count within this namespace, and can be used to understand where clients are originating. Authentication methods are organized by path.',
          newCopy: `The new clients used by the auth method for this ${dateText}. This aids in understanding which auth methods create and use new clients${
            dateText === 'date range' ? ' over time.' : '.'
          }`,
          totalCopy: `The total clients used by the auth method for this ${dateText}. This number is useful for identifying overall usage volume. `,
        };
      case false:
        return {
          description:
            'This data shows the top ten namespaces by client count and can be used to understand where clients are originating. Namespaces are identified by path. To see all namespaces, export this data.',
          newCopy: `The new clients in the namespace for this ${dateText}.
          This aids in understanding which namespaces create and use new clients${
            dateText === 'date range' ? ' over time.' : '.'
          }`,
          totalCopy: `The total clients in the namespace for this ${dateText}. This number is useful for identifying overall usage volume.`,
        };
      default:
        return '';
    }
  }
}
