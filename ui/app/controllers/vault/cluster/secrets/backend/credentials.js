import Controller from '@ember/controller';

export default Controller.extend({
  queryParams: ['action'],
  action: '',
  reset() {
    this.set('action', '');
  },
});
