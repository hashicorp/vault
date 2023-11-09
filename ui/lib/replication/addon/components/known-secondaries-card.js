/**
 * Copyright (c) HashiCorp, Inc.
 * SPDX-License-Identifier: BUSL-1.1
 */

import Component from '@ember/component';

/**
 * @module KnownSecondariesCard
 * KnownSecondariesCard components are used on the Replication Details dashboards to display a table of known secondary clusters.
 *
 * @example
 * <KnownSecondariesCard @cluster={{clusterModel}} @replicationAttrs={{replicationAttrs}} />
 *
 * @param {object} cluster=null - The cluster model.
 * @param {object} replicationAttrs=null - The attributes passed directly from the cluster model. These are passed down to the KnownSecondariesTable.
 */

export default Component.extend({
  tagName: '',
  cluster: null,
  replicationAttrs: null,
});
