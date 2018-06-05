import Ember from 'ember';
const { assert, inject, Component } = Ember;

export default Component.extend({
  tagName: '',
  flashMessages: inject.service(),
  params: null,
  successMessage() {
    return 'Save was successful';
  },
  errorMessage() {
    return 'There was an error saving';
  },
  onError(model) {
    if (model && model.rollbackAttributes) {
      model.rollbackAttributes();
    }
  },
  onSuccess() {},
  // override and return a promise
  transaction() {
    assert('override transaction call in an extension of popup-base', false);
  },

  actions: {
    performTransaction() {
      let args = [...arguments];
      let messageArgs = this.messageArgs(...args);
      return this.transaction(...args)
        .then(() => {
          this.get('onSuccess')();
          this.get('flashMessages').success(this.successMessage(...messageArgs));
        })
        .catch(e => {
          this.onError(...messageArgs);
          this.get('flashMessages').success(this.errorMessage(e, ...messageArgs));
        });
    },
  },
});
