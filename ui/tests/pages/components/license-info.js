/**
 * Copyright IBM Corp. 2016, 2025
 * SPDX-License-Identifier: BUSL-1.1
 */

import { text, collection } from 'ember-cli-page-object';

export default {
  detailRows: collection('[data-test-detail-row]', {
    rowName: text('[data-test-row-label]'),
    rowValue: text('.column.is-flex-center'),
  }),
  featureRows: collection('[data-test-feature-row]', {
    featureName: text('[data-test-row-label]'),
    featureStatus: text('[data-test-feature-status]'),
  }),
};
