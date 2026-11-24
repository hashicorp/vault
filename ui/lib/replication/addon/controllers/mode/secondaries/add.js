/**
 * Copyright IBM Corp. 2016, 2026
 * SPDX-License-Identifier: BUSL-1.1
 */

import ReplicationController from 'replication/controllers/application';
import { isDurationOverMax } from 'core/utils/duration-utils';

// Vault's default token TTL, sent in place of a value that was entered, turned off and is too large to send.
const DEFAULT_TTL = '1800s';

export default ReplicationController.extend({
  isTtlValid: true,

  actions: {
    updateTtl: function (ttl) {
      // A TTL that is turned off keeps being sent as before, unless Vault would reject it.
      const isHiddenOverMax = !ttl.enabled && isDurationOverMax(ttl.seconds);
      this.setProperties({
        ttl: isHiddenOverMax ? DEFAULT_TTL : `${ttl.seconds}s`,
        isTtlValid: ttl.isValid,
      });
    },

    onSubmit(action, mode, data, event) {
      // the picker explains an over-limit TTL, so it only needs to stop the request
      if (!this.isTtlValid) {
        event?.preventDefault?.();
        return;
      }
      return this._super(...arguments);
    },
  },
});
