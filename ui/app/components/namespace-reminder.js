import { inject as service } from '@ember/service';
import { not } from '@ember/object/computed';
import Component from '@ember/component';
import { computed } from '@ember/object';

export default Component.extend({
  namespace: service(),
  showMessage: not('namespace.inRootNamespace'),
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
