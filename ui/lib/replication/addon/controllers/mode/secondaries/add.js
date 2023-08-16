import ReplicationController from 'replication/controllers/application';

export default ReplicationController.extend({
  actions: {
    updateTtl: function (ttl) {
      this.set('ttl', `${ttl.seconds}s`);
    },
  },
});
