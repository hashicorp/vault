import { inject as service } from '@ember/service';
import { alias } from '@ember/object/computed';
import Component from '@ember/component';
import { computed } from '@ember/object';

export default Component.extend({
  namespaceService: service('namespace'),
  currentNamespace: alias('namespaceService.path'),

  tagName: '',
  //public api
  targetNamespace: null,
  showLastSegment: false,

  normalizedNamespace: computed('targetNamespace', function() {
    let ns = this.get('targetNamespace');
    return (ns || '').replace(/\.+/g, '/').replace(/☃/g, '.');
  }),

  namespaceDisplay: computed('normalizedNamespace', 'showLastSegment', function() {
    let ns = this.get('normalizedNamespace');
    if (!ns) return 'root';
    let showLastSegment = this.get('showLastSegment');
    let parts = ns?.split('/');
    return showLastSegment ? parts[parts.length - 1] : ns;
  }),

  isCurrentNamespace: computed('targetNamespace', 'currentNamespace', function() {
    return this.get('currentNamespace') === this.get('targetNamespace');
  }),
});
