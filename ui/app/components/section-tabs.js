import Ember from 'ember';

const SectionTabs = Ember.Component.extend({
  tagName: '',

  model: null,
  tabType: 'authSettings',
});

SectionTabs.reopenClass({
  positionalParams: ['model', 'tabType'],
});

export default SectionTabs;
