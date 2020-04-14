/**
 * @module ReplicationHeader
 * ARG TODO: finish
 *
 * @example
 * ```js
 * <ReplicationHeader finish/>
 * ```
 */

import Component from '@ember/component';
import { computed } from '@ember/object';
import layout from '../templates/components/replication-header';

export default Component.extend({
  layout,
  data: null,
  isSecondary: computed('data', function() {
    let data = this.data;
    if (data.dr.mode === 'secondary' && data.rm.mode == 'dr') {
      return true;
    }
  }),
});
