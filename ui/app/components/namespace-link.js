import Ember from 'ember';

const { Component, computed, inject } = Ember;

export default Component.extend({
  namespaceService: inject.service('namespace'),
  currentNamespace: computed.alias('namespaceService.path'),

  tagName: '',
  //public api
  targetNamespace: null,
  showLastSegment: false,

  normalizedNamespace: computed('targetNamespace', function() {
    let ns = this.get('targetNamespace');
    return (ns || '').replace(/\.+/g, '/').replace('☃', '.');
  }),

  namespaceDisplay: computed('normalizedNamespace', 'showLastSegment', function() {
    let ns = this.get('normalizedNamespace');
    let showLastSegment = this.get('showLastSegment');
    let parts = ns.split('/');
    if (ns === '') {
      return 'root';
    }
    return showLastSegment ? parts[parts.length - 1] : ns;
  }),

  isCurrentNamespace: computed('targetNamespace', 'currentNamespace', function() {
    return this.get('currentNamespace') === this.get('targetNamespace');
  }),
});
