import Ember from 'ember';

export default Ember.Component.extend({
  tagName: '',

  // api
  isCertTab: false,
  isConfigure: false,
  baseKey: null,
  backendCrumb: null,
  model: null,
});
