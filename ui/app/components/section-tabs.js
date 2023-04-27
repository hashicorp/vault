/**
 * Copyright (c) HashiCorp, Inc.
 * SPDX-License-Identifier: MPL-2.0
 */

import Component from '@ember/component';

const SectionTabs = Component.extend({
  tagName: '',
  model: null,
  tabType: 'authSettings',
});

SectionTabs.reopenClass({
  positionalParams: ['model', 'tabType', 'paths'],
});

export default SectionTabs;
