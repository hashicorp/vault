/**
 * Copyright IBM Corp. 2016, 2025
 * SPDX-License-Identifier: BUSL-1.1
 */

import DateBase from './-date-base';

export default DateBase.extend({
  compute() {
    this._super(...arguments);

    return Date.now();
  },
});
