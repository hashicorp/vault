import Component from '@ember/component';
import { computed } from '@ember/object';

/**
 * @module KnownSecondariesTable
 * KnownSecondariesTable components are used on the Replication Details dashboards to display a table of known secondary clusters.
 *
 * @example
 * ```js
 * <KnownSecondariesTable @replicationAttrs={{replicationAttrs}} />
 * ```
 * @param {object} replicationAttrs=null - The attributes passed directly from the cluster model used to access the array of known secondaries.
 */

export default Component.extend({
  replicationAttrs: null,
  knownSecondaries: computed('replicationAttrs', function() {
    const { replicationAttrs } = this;
    return replicationAttrs.knownSecondaries;
  }),
});
