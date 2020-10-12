import Controller from '@ember/controller';

export default Controller.extend({
  queryParams: ['start', 'end'],

  start: null,
  end: null,
});
