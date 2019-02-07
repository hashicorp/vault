import AuthConfigComponent from './config';
import DS from 'ember-data';

export default AuthConfigComponent.extend({
  tuneModel() {},
  actions: {
    tuneAndRedirect: function() {
      debugger; // eslint-disable-line
      return this.tuneModel(this.model)
        .then(() => {
          this.flashMessages.success('The configuration options were saved successfully.');
        })
        .catch(err => {
          // AdapterErrors are handled by the error-message component
          // in the form
          if (err instanceof DS.AdapterError === false) {
            throw err;
          }
          return;
        });
    },
  },
});
