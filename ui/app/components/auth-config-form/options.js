import AuthConfigComponent from './config';
import { task } from 'ember-concurrency';
import DS from 'ember-data';

export default AuthConfigComponent.extend({
  saveModel: task(function*() {
    const model = this.get('model');
    let data = model.get('config').serialize();
    data.description = model.get('description');
    try {
      yield model.tune(data);
    } catch (err) {
      // AdapterErrors are handled by the error-message component
      // in the form
      if (err instanceof DS.AdapterError === false) {
        throw err;
      }
      return;
    }
    this.get('flashMessages').success('The configuration options were saved successfully.');
  }),
});
