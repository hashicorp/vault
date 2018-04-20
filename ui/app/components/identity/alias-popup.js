import Ember from 'ember';

export default Ember.Component.extend({
  flashMessages: Ember.inject.service(),
  item: null,

  tagName: '',
  actions: {
    delete(item) {
      let type = item.get('identityType');
      let id = item.id;
      return item
        .destroyRecord()
        .then(() => {
          this.get('flashMessages').success(`Successfully deleted ${type}: ${id}`);
        })
        .catch(e => {
          this.get('flashMessages').success(
            `There was a problem deleting ${type}: ${id} - ${e.error.join(' ') || e.message}`
          );
        });
    },
  },
});
