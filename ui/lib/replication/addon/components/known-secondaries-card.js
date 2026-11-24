/**
 * Copyright IBM Corp. 2016, 2025
 * SPDX-License-Identifier: BUSL-1.1
 */

import Component from '@ember/component';

/**
 * @module KnownSecondariesCard
 * KnownSecondariesCard components are used on the Replication Details dashboards to display a table of known secondary clusters.
 *
 * @example
 * ```js
 * <KnownSecondariesCard @cluster={{cluster}} @replicationAttrs={{replicationAttrs}} />
 * ```
 * @param {object} cluster=null - The cluster resource.
 * @param {object} replicationAttrs=null - The replication attributes passed directly from the cluster resource. These are passed down to the KnownSecondariesTable.
 */

export default Component.extend({
  tagName: '',
  cluster: null,
  replicationAttrs: null,
});
