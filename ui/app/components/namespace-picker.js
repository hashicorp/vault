import Ember from 'ember';

const { Component, inject } = Ember;

export default Component.extend({
  namespaceService: inject.service(),
  //passed from the queryParam
  namespace: null,

  didRecieveAttrs() {
    this._super(...arguments);
    this.get('namespaceService').setNamespace(this.get('namespace'));
  },
});
