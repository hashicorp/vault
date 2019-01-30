import Controller from '@ember/controller';
import ListController from 'vault/mixins/list-controller';

export default Controller.extend(ListController, {
  actions: {
    onDelete() {
      this.send('reload');
    },
  },
});
