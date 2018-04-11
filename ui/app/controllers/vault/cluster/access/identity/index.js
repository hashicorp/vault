import Ember from 'ember';
import ListController from 'vault/mixins/list-controller';

const { inject } = Ember;

export default Ember.Controller.extend(ListController, {
  flashMessages: inject.service(),

  actions: {
    delete(model) {
      let type = model.get('identityType');
      let id = model.id;
      return model
        .destroyRecord()
        .then(() => {
          this.send('willTransition');
          this.get('flashMessages').success(`Successfully deleted ${type}: ${id}`);
        })
        .catch(e => {
          this.get('flashMessages').success(
            `There was a problem deleting ${type}: ${id} - ${e.error.join(' ') || e.message}`
          );
        });
    },
    reloadRecord(model) {
      model.reload();
    },
  },
});
