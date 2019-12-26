import { computed } from '@ember/object';
import Controller from '@ember/controller';

export default Controller.extend({
  isConfigurable: computed('model.type', function() {
    const configurableEngines = ['aws', 'ssh', 'pki'];
    return configurableEngines.includes(this.get('model.type'));
  }),
});
