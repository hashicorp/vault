import Ember from 'ember';

const { Component, inject, computed } = Ember;

export default Component.extend({
  namespace: inject.service(),
  showMessage: computed.not('namespace.inRootNamespace'),
});
