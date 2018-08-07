import Ember from 'ember';
import ListController from 'vault/mixins/list-controller';

export default Ember.Controller.extend(ListController, {
  actions: {
    onDelete() {
      this.send('reload');
    },
  },
});
