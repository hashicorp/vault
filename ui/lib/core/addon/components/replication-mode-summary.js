/**
 * Copyright (c) HashiCorp, Inc.
 * SPDX-License-Identifier: BUSL-1.1
 */

import { service } from '@ember/service';
import { equal } from '@ember/object/computed';
import { get, computed } from '@ember/object';
import Component from '@ember/component';
import layout from '../templates/components/replication-mode-summary';

const replicationAttr = function (attr) {
  return computed(`cluster.{dr,performance}.${attr}`, 'cluster', 'mode', function () {
    const { mode, cluster } = this;
    return get(cluster, `${mode}.${attr}`);
  });
};
export default Component.extend({
  layout,
  version: service(),
  router: service(),
  namespace: service(),
  classNameBindings: ['isMenu::box'],
  attributeBindings: ['href', 'target'],
  display: 'banner',
  isMenu: equal('display', 'menu'),
  href: computed(
    'cluster.id',
    'display',
    'mode',
    'replicationEnabled',
    'version.hasPerfReplication',
    function () {
      const display = this.display;
      const mode = this.mode;
      if (mode === 'performance' && display === 'menu' && this.version.hasPerfReplication === false) {
        return 'https://www.hashicorp.com/products/vault';
      }
      if (this.replicationEnabled || display === 'menu') {
        return this.router.urlFor('vault.cluster.replication.mode.index', this.cluster.id, mode);
      }
      return null;
    }
  ),
  target: computed('isPerformance', 'version.hasPerfReplication', function () {
    if (this.isPerformance && this.version.hasPerfReplication === false) {
      return '_blank';
    }
    return null;
  }),
  internalLink: false,
  isPerformance: equal('mode', 'performance'),
  replicationEnabled: replicationAttr('replicationEnabled'),
  replicationUnsupported: equal('cluster.mode', 'unsupported'),
  replicationDisabled: replicationAttr('replicationDisabled'),
  syncProgressPercent: replicationAttr('syncProgressPercent'),
  syncProgress: replicationAttr('syncProgress'),
  secondaryId: replicationAttr('secondaryId'),
  modeForUrl: replicationAttr('modeForUrl'),
  clusterIdDisplay: replicationAttr('clusterIdDisplay'),
  mode: null,
  cluster: null,
  modeState: computed('cluster', 'mode', function () {
    const { cluster, mode } = this;
    const clusterState = cluster[mode].state;
    return clusterState;
  }),
});
