import Component from '@ember/component';

const SectionTabs = Component.extend({
  tagName: '',

  model: null,
  tabType: 'authSettings',
});

SectionTabs.reopenClass({
  positionalParams: ['model', 'tabType'],
});

export default SectionTabs;
