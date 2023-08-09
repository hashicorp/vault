/**
 * Copyright (c) HashiCorp, Inc.
 * SPDX-License-Identifier: MPL-2.0
 */

import Component from '@ember/component';
import layout from '../templates/components/replication-header';

/**
 * @module ReplicationHeader
 * The `ReplicationHeader` is a header component used on the Replication Dashboards.
 *
 * @example
 * ```js
 * <ReplicationHeader
    @data={{model}}
    @title="Secondary"
    @secondaryID="meep_123"
    @isSummaryDashboard=false
    />
 * ```
 * @param {Object} model=null - An Ember data object pulled from the Ember cluster model.
 * @param {String} title=null - The title of the header.
 * @param {String} [secondaryID=null] - The secondaryID pulled off of the model object. 
 * @param {Boolean} isSummaryDashboard=false - True when you have both a primary performance and dr cluster dashboard.
 */

export default Component.extend({
  layout,
  data: null,
  classNames: ['replication-header'],
  isSecondary: null,
  secondaryId: null,
  isSummaryDashboard: false,
  'data-test-replication-header': true,
  attributeBindings: ['data-test-replication-header'],
});
