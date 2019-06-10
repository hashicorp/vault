import DS from 'ember-data';

export default DS.Model.extend({
  getHelpUrl: function(backend) {
    return `/v1/${backend}/scope/example/role/example/credentials?help=1`;
  },
});
