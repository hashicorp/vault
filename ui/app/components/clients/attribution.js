/**
 * Copyright (c) HashiCorp, Inc.
 * SPDX-License-Identifier: BUSL-1.1
 */

import Component from '@glimmer/component';

/**
 * @module Attribution
 * Attribution components display the top 10 total client counts for namespaces or mounts during a billing period.
 * A horizontal bar chart shows on the right, with the top namespace/auth method and respective client totals on the left.
 *
 * @example
 *  <Clients::Attribution
 *    @noun="auth mount"
 *    @attribution={{array (hash label="ns1" clients=100)}}
 *    @responseTimestamp="2018-04-03T14:15:30"
 *    @isSecretsSyncActivated={{true}}
 *  />
 *
 * @param {string} noun - noun which reflects the type of data and used in title. Should be "namespace" (default) or "auth mount"
 * @param {array} attribution - array of objects containing a label and breakdown of client counts for total clients
 * @param {string} responseTimestamp -  ISO timestamp created in serializer to timestamp the response, renders in bottom left corner below attribution chart
 * @param {boolean} isSecretsSyncActivated - boolean reflecting if secrets sync is activated. Determines the labels and data shown
 */

export default class Attribution extends Component {
  get noun() {
    return this.args.noun || 'namespace';
  }

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

  get sortedAttribution() {
    return this.args.attribution?.sort((a, b) => b.clients - a.clients);
  }

  // truncate data before sending to chart component
  get topTenAttribution() {
    return this.sortedAttribution?.slice(0, 10);
  }

  get topAttribution() {
    // get top namespace or auth method
    return this.sortedAttribution ? this.sortedAttribution[0] : null;
  }

  get chartText() {
    if (this.noun === 'namespace') {
      return {
        subtext: 'This data shows the top ten namespaces by total clients for the date range selected.',
        description:
          'This data shows the top ten namespaces by total clients and can be used to understand where clients are originating. Namespaces are identified by path. To see all namespaces, click export data on top of the page.',
      };
    } else {
      return {
        subtext:
          'The total clients used by the auth method for this date range. This number is useful for identifying overall usage volume.',
        description:
          'This data shows the top ten authentication methods by client count within this namespace, and can be used to understand where clients are originating. Auth methods are organized by path. To see all auth methods, click export data on top of the page.',
      };
    }
  }
}
