import Controller from '@ember/controller';

export default Controller.extend({
  queryParams: ['state', 'code'],
  code: null,
  state: null,
});
