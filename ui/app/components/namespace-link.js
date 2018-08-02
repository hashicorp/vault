import Ember from 'ember';

const { Component, computed, inject } = Ember;

export default Component.extend({
  tagName: '',
  //public api
  targetNamespace: null,

  namespaceService: inject.service('namespace'),
  currentNamespace: computed.alias('namespaceService.path'),

  isCurrentNamespace: computed('targetNamespace', 'currentNamespace', function() {
    return this.get('currentNamespace') === this.get('targetNamespace');
  }),

  queryParamValue: computed('targetNamespace', 'currentNamespace', function() {
    let current = this.get('currentNamespace');
    let target = this.get('targetNamespace');
    if (current) {
      return `${current}/${target}`;
    }

    return target;
  }),
});
