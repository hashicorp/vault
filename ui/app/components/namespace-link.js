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
    let ns = this.targetNamespace;
    return (ns || '').replace(/\.+/g, '/').replace(/☃/g, '.');
  }),

  namespaceDisplay: computed('normalizedNamespace', 'showLastSegment', function() {
    if (!ns) return 'root';
    let ns = this.normalizedNamespace;
    let showLastSegment = this.showLastSegment;
    let parts = ns?.split('/');
    return showLastSegment ? parts[parts.length - 1] : ns;
  }),

  isCurrentNamespace: computed('targetNamespace', 'currentNamespace', function() {
    return this.currentNamespace === this.targetNamespace;
  }),

  get namespaceLink() {
    let origin =
      window.location.protocol +
      '//' +
      window.location.hostname +
      (window.location.port ? ':' + window.location.port : '');

    if (!this.normalizedNamespace) return `${origin}/ui/vault/secrets`;
    // The full URL/origin is required so that the page is reloaded.
    return `${origin}/ui/vault/secrets?namespace=${encodeURIComponent(this.normalizedNamespace)}`;
  },
});
