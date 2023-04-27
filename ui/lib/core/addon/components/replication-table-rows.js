/**
 * Copyright (c) HashiCorp, Inc.
 * SPDX-License-Identifier: MPL-2.0
 */

import { reads } from '@ember/object/computed';
import Component from '@ember/component';
import { computed } from '@ember/object';
import layout from '../templates/components/replication-table-rows';

/**
 * @module ReplicationTableRows
 * The `ReplicationTableRows` component is table component.  It displays cluster mode details specific to the cluster of the Dashboard it is used on.
 *
 * @example
 * ```js
 * <ReplicationTableRows
    @replicationDetails={{replicationDetails}}
    @clusterMode="primary"
    />
 * ```
 * @param {Object} replicationDetails=null - An Ember data object pulled from the Ember Model. It contains details specific to the whether the replication is dr or performance.
 * @param {String} clusterMode=null - The cluster mode (e.g. primary or secondary) passed through to a table component.
 */

export default Component.extend({
  layout,
  classNames: ['replication-table-rows'],
  replicationDetails: null,
  clusterMode: null,
  secondaryId: reads('replicationDetails.secondaryId'),
  primaryClusterAddr: computed('replicationDetails.primaryClusterAddr', function () {
    return this.replicationDetails.primaryClusterAddr || 'None set';
  }),
  merkleRoot: computed('replicationDetails.merkleRoot', function () {
    return this.replicationDetails.merkleRoot || 'unknown';
  }),
  clusterId: computed('replicationDetails.clusterId', function () {
    return this.replicationDetails.clusterId || 'unknown';
  }),
});
