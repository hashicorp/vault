import Component from '@ember/component';
import layout from '../templates/components/info-table';

/**
 * @module InfoTable
 * InfoTable components are a table with a single column and header. They are used to render a list of InfoTableRow components.
 *
 * @example
 * ```js
 * <InfoTable @replicationAttrs={{replicationAttrs}} />
 * ```
 * @param {object} replicationAttrs=null - The attributes passed directly from the cluster model used to access the array of known secondaries. We use this to grab the secondaries.
 */

export default Component.extend({
  layout,
  tagName: '',
  title: null,
  header: null,
  items: null,
});
