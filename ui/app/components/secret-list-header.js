import Component from '@ember/component';

export default Component.extend({
  tagName: '',

  // api
  isCertTab: false,
  isConfigure: false,
  baseKey: null,
  backendCrumb: null,
  model: null,
});
