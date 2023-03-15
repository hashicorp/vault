/**
 * Copyright (c) HashiCorp, Inc.
 * SPDX-License-Identifier: MPL-2.0
 */

import Component from '@ember/component';
import layout from '../templates/components/info-table';

/**
 * @module InfoTable
 * InfoTable components are a table with a single column and header. They are used to render a list of InfoTableRow components.
 *
 * @example
 * ```js
 * <InfoTable
        @title="Known Primary Cluster Addrs"
        @header="cluster_addr"
        @items={{knownPrimaryClusterAddrs}}
      />
 * ```
 * @param {String} [title=Info Table] - The title of the table. Used for accessibility purposes.
 * @param {String} header=null - The column header.
 * @param {Array} items=null - An array of strings which will be used as the InfoTableRow value.
 */

export default Component.extend({
  layout,
  tagName: '',
  title: 'Info Table',
  header: null,
  items: null,
});
