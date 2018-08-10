import Ember from 'ember';

const { Component, computed, inject } = Ember;

export default Component.extend({
  tagName: '',
  //public api
  targetNamespace: null,

  normalizedNamespace: computed('targetNamespace', function() {
    let ns = this.get('targetNamespace');
    return (ns || '').replace(/\.+/g, '/').replace('☃', '.');
  }),
  namespaceService: inject.service('namespace'),
  currentNamespace: computed.alias('namespaceService.path'),

  isCurrentNamespace: computed('targetNamespace', 'currentNamespace', function() {
    return this.get('currentNamespace') === this.get('targetNamespace');
  }),
});
