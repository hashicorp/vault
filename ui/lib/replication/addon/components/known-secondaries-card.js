import Component from '@ember/component';

/**
 * @module KnownSecondariesCard
 * KnownSecondariesCard components
 *
 * @example
 * ```js
 * <KnownSecondariesCard @replicationAttrs={{replicationAttrs}} />
 * ```
 * @param {string} [replicationAttrs=null] - The attributes passed directly from the cluster model. These are passed down to the KnownSecondariesTable.
 */

export default Component.extend({
  tagName: '',
  replicationAttrs: null,
});
