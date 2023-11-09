/**
 * Copyright (c) HashiCorp, Inc.
 * SPDX-License-Identifier: BUSL-1.1
 */

import Component from '@ember/component';
import { computed } from '@ember/object';
import { clusterStates } from 'core/helpers/cluster-states';

/**
 * @module ReplicationPrimaryCard
 * The `ReplicationPrimaryCard` component is a card-like component.  It displays cluster mode details specific for DR and Performance Primaries.
 *
 * @example
 * <ReplicationPrimaryCard
    @title='Last WAL entry'
    @description='Index of last Write Ahead Logs entry written on local storage.'
    @metric={{replicationAttrs.lastWAL}}
    />
 *
 * @param {String} [title=null] - The title to be displayed on the top left corner of the card.
 * @param {String} [description=null] - Helper text to describe the metric on the card.
 * @param {String} [glyph=null] - The glyph to display beside the metric.
 * @param {String} metric=null - The main metric to highlight on the card.
 */

export default Component.extend({
  tagName: '',
  title: null,
  description: null,
  metric: null,
  glyph: null,
  hasError: computed('title', 'metric', function () {
    const { title, metric } = this;

    // only show errors on the state card if state is not ok
    if (title === 'State') {
      return metric && !clusterStates([metric]).isOk;
    }
    return false;
  }),
  errorMessage: computed('hasError', function () {
    return this.hasError ? 'Check server logs!' : false;
  }),
});
