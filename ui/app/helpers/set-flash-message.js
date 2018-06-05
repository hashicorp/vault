import Ember from 'ember';

const { Helper, inject } = Ember;

export default Helper.extend({
  flashMessages: inject.service(),

  compute([message, type]) {
    return () => {
      this.get('flashMessages')[type || 'success'](message);
    };
  },
});
