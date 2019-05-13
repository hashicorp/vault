import Controller from '@ember/controller';

export default Controller.extend({
  actions: {
    onSave({ saveType }) {
      if (saveType === 'destroyRecord') {
        this.send('reload');
      }
    },
  },
});
