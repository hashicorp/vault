import { inject as service } from '@ember/service';
import Component from '@ember/component';
import DS from 'ember-data';

const AuthConfigBase = Component.extend({
  tagName: '',
  model: null,

  flashMessages: service(),
  wizard: service(),
  saveModel() {},
  actions: {
    saveAndRedirect: function() {
      return this.saveModel(this.model)
        .then(() => {
          this.flashMessages.success('The configuration was saved successfully.');
        })
        .catch(err => {
          if (err instanceof DS.AdapterError === false) {
            throw err;
          }
          return;
        });
    },
  },
});

AuthConfigBase.reopenClass({
  positionalParams: ['model'],
});

export default AuthConfigBase;
