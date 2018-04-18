import Ember from 'ember';

export default Ember.Component.extend({
  model: null,
  errors: [],
  errorMessage: null,

  displayErrors: Ember.computed(
    'errorMessage',
    'model.isError',
    'model.adapterError.errors.@each',
    'errors',
    'errors.@each',
    function() {
      const errorMessage = this.get('errorMessage');
      const errors = this.get('errors');
      const modelIsError = this.get('model.isError');
      if (errorMessage) {
        return [errorMessage];
      }

      if (errors && errors.length > 0) {
        return errors;
      }

      if (modelIsError) {
        return this.get('model.adapterError.errors');
      }
    }
  ),
});
