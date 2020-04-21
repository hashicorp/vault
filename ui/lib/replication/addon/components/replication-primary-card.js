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
 * @param {object} [glyph=null] - The glyph to display beside the metric.
 * @param {string} metric=null - The main metric to highlight on the card.
 */

export default Component.extend({
  tagName: '',
  title: null,
  description: null,
  metric: null,
  glyph: null,
  hasError: computed('title', 'metric', function() {
    // TODO: can we make a map somewhere in the cluster that keeps track of all the good and bad states
    // as well as their glyphs? this could replace the cluster StateDiplay and StateGlyph
    // TODO: then add tests to ensure we show the correct error msg
    return this.title === 'State' && this.metric !== 'running';
  }),
  errorMessage: computed('hasError', function() {
    // TODO figure out if we need another error message
    return this.hasError ? 'Check server logs!' : false;
  }),
});
