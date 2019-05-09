import Component from '@ember/component';
import { computed } from '@ember/object';

export default Component.extend({
  tagName: '',

  // api
  isCertTab: false,
  isConfigure: false,
  baseKey: null,
  backendCrumb: null,
  model: null,
  options: null,
  hasItems: computed('model.meta.total', function() {
    return this.get('model.meta.total');
  }),
  isConfigurable: computed('model.type', function() {
    const configurableEngines = ['aws', 'ssh', 'pki'];
    return configurableEngines.includes(this.get('model.type'));
  }),
  isConfigurableTab: computed('isCertTab', 'isConfigure', function() {
    return this.get('isCertTab') || this.get('isConfigure');
  }),
});
