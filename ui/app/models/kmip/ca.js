/**
 * Copyright (c) HashiCorp, Inc.
 * SPDX-License-Identifier: BUSL-1.1
 */

import Model, { belongsTo, attr } from '@ember-data/model';

export default Model.extend({
  config: belongsTo('kmip/config', { async: false, inverse: 'ca' }),
  caPem: attr('string', {
    label: 'CA PEM',
  }),
});
