import DS from 'ember-data';

export default DS.Model.extend({
  useOpenAPI: true,
  getHelpUrl(path) {
    return `/v1/${path}/config?help=1`;
  },
});
