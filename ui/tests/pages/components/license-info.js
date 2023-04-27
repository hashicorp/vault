/**
 * Copyright (c) HashiCorp, Inc.
 * SPDX-License-Identifier: MPL-2.0
 */

import { text, collection } from 'ember-cli-page-object';

export default {
  detailRows: collection('[data-test-detail-row]', {
    rowName: text('[data-test-row-label]'),
    rowValue: text('.column.is-flex'),
  }),
  featureRows: collection('[data-test-feature-row]', {
    featureName: text('[data-test-row-label]'),
    featureStatus: text('[data-test-feature-status]'),
  }),
};
