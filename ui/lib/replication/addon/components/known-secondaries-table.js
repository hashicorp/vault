/**
 * Copyright (c) HashiCorp, Inc.
 * SPDX-License-Identifier: BUSL-1.1
 */

import Component from '@ember/component';

/**
 * @module KnownSecondariesTable
 * KnownSecondariesTable components are used on the Replication Details dashboards
 * to display a table of known secondary clusters.
 *
 * @example
 * <KnownSecondariesTable @replicationAttrs={{replicationAttrs}} />
 *
 * @param {array} secondaries=null - The array of secondaries from the replication
 * status endpoint. Contains the secondary api_address, id and connected_state.
 */

export default Component.extend({
  secondaries: null,
});
