/**
 * Copyright (c) HashiCorp, Inc.
 * SPDX-License-Identifier: MPL-2.0
 */

import DateBase from './-date-base';

export default DateBase.extend({
  compute() {
    this._super(...arguments);

    return Date.now();
  },
});
