/**
 * Copyright (c) HashiCorp, Inc.
 * SPDX-License-Identifier: MPL-2.0
 */

import Component from '@ember/component';

/**
 * @module KnownSecondariesTable
 * KnownSecondariesTable components are used on the Replication Details dashboards
 * to display a table of known secondary clusters.
 *
 * @example
 * ```js
 * <KnownSecondariesTable @replicationAttrs={{replicationAttrs}} />
 * ```
 * @param {array} secondaries=null - The array of secondaries from the replication
 * status endpoint. Contains the secondary api_address, id and connected_state.
 */

export default Component.extend({
  secondaries: null,
});
