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
    let dr = this.data.dr;
    if (!dr) {
      return false;
    }
    if (dr.mode === 'secondary' && data.rm.mode == 'dr') {
      return true;
    }
  }),
});
