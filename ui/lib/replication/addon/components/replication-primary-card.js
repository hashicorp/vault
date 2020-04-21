import Component from '@ember/component';

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
});
