import { inject as service } from '@ember/service';
import { equal } from '@ember/object/computed';
import { getProperties, get, computed } from '@ember/object';
import Component from '@ember/component';

const replicationAttr = function(attr) {
  return computed('mode', `cluster.{dr,performance}.${attr}`, function() {
    const { mode, cluster } = getProperties(this, 'mode', 'cluster');
    return get(cluster, `${mode}.${attr}`);
  });
};
export default Component.extend({
  version: service(),
  router: service(),
  namespace: service(),
  classNameBindings: ['isMenu::box', 'isMenu::level'],
  attributeBindings: ['href', 'target'],
  display: 'banner',
  isMenu: equal('display', 'menu'),
  href: computed('display', 'mode', 'replicationEnabled', 'version.hasPerfReplication', function() {
    const display = this.get('display');
    const mode = this.get('mode');
    if (mode === 'performance' && display === 'menu' && this.get('version.hasPerfReplication') === false) {
      return 'https://www.hashicorp.com/products/vault';
    }
    if (this.get('replicationEnabled') || display === 'menu') {
      return this.get('router').urlFor(
        'vault.cluster.replication.mode.index',
        this.get('cluster.name'),
        mode
      );
    }
    return null;
  }),
  target: computed('isPerformance', 'version.hasPerfReplication', function() {
    if (this.get('isPerformance') && this.get('version.hasPerfReplication') === false) {
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
  partialName: computed('display', function() {
    return this.get('display') === 'menu'
      ? 'partials/replication/replication-mode-summary-menu'
      : 'partials/replication/replication-mode-summary';
  }),
});
