import Ember from 'ember';
import { hrefTo } from 'vault/helpers/href-to';
const { computed, get, getProperties } = Ember;

const replicationAttr = function(attr) {
  return computed('mode', `cluster.{dr,performance}.${attr}`, function() {
    const { mode, cluster } = getProperties(this, 'mode', 'cluster');
    return get(cluster, `${mode}.${attr}`);
  });
};
export default Ember.Component.extend({
  version: Ember.inject.service(),
  classNames: ['level', 'box-label'],
  classNameBindings: ['isMenu:is-mobile'],
  attributeBindings: ['href', 'target'],
  display: 'banner',
  isMenu: computed.equal('display', 'menu'),
  href: computed('display', 'mode', 'replicationEnabled', 'version.hasPerfReplication', function() {
    const display = this.get('display');
    const mode = this.get('mode');
    if (mode === 'performance' && display === 'menu' && this.get('version.hasPerfReplication') === false) {
      return 'https://www.hashicorp.com/products/vault';
    }
    if (this.get('replicationEnabled') || display === 'menu') {
      return hrefTo(this, 'vault.cluster.replication.mode.index', this.get('cluster.name'), mode);
    }
    return null;
  }),
  target: computed('isPerformance', 'version.hasPerfReplication', function() {
    if (this.get('isPerformance') && this.get('version.hasPerfReplication') === false) {
      return '_blank';
    }
    return null;
  }),
  isPerformance: computed.equal('mode', 'performance'),
  replicationEnabled: replicationAttr('replicationEnabled'),
  replicationUnsupported: computed.equal('cluster.mode', 'unsupported'),
  replicationDisabled: replicationAttr('replicationDisabled'),
  syncProgressPercent: replicationAttr('syncProgressPercent'),
  syncProgress: replicationAttr('syncProgress'),
  secondaryId: replicationAttr('secondaryId'),
  modeForUrl: replicationAttr('modeForUrl'),
  clusterIdDisplay: replicationAttr('clusterIdDisplay'),
  mode: null,
  cluster: null,
  partialName: computed('display', function() {
    return this.get('display') === 'menu'
      ? 'partials/replication/replication-mode-summary-menu'
      : 'partials/replication/replication-mode-summary';
  }),
});
