/**
 * Copyright (c) HashiCorp, Inc.
 * SPDX-License-Identifier: MPL-2.0
 */

import Model, { belongsTo, attr } from '@ember-data/model';

export default Model.extend({
  config: belongsTo('kmip/config', { async: false }),
  caPem: attr('string', {
    label: 'CA PEM',
  }),
});
