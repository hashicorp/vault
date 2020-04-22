import Component from '@ember/component';
import { computed } from '@ember/object';

/**
 * @module ReplicationPrimaryCard
 * ReplicationPrimaryCard components
 *
 * @example
 * ```js
 * <ReplicationPrimaryCard
    @title='Last WAL entry'
    @description='Index of last Write Ahead Logs entry written on local storage.'
    @metric={{replicationAttrs.lastWAL}}
    />
 * ```
 * @param {string} [title=null] - The title to be displayed on the top left corner of the card.
 * @param {string} [description=null] - Helper text to describe the metric on the card.
 * @param {string} [glyph=null] - The glyph to display beside the metric.
 * @param {boolean} [hasOkState=true] - Whether or not the cluster is in an okay state. This comes from the cluster model.
 * @param {string} metric=null - The main metric to highlight on the card.
 */

export default Component.extend({
  tagName: '',
  title: null,
  description: null,
  metric: null,
  glyph: null,
  hasOkState: true,
  hasError: computed('hasOkState', 'title', 'metric', function() {
    return this.title === 'State' && !this.hasOkState;
  }),
  errorMessage: computed('hasError', function() {
    // TODO figure out if we need another error message
    return this.hasError ? 'Check server logs!' : false;
  }),
});
