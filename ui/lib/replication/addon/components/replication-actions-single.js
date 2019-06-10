import Component from '@ember/component';

export default Component.extend({
  onSubmit() {},
  replicationMode: null,
  replicationDisplayMode: null,
  model: null,

  actions: {
    onSubmit() {
      return this.onSubmit(...arguments);
    },
  },
});
