import Ember from 'ember';

const { Component, inject, computed } = Ember;

export default Component.extend({
  namespace: inject.service(),
  showMessage: computed.not('namespace.inRootNamespace'),
  //public API
  noun: null,
  mode: 'edit',
  modeVerb: computed(function() {
    let mode = this.get('mode');
    if (!mode) {
      return '';
    }
    return mode.endsWith('e') ? `${mode}d` : `${mode}ed`;
  }),
});
