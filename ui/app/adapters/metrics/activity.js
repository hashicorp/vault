import Application from '../application';

export default Application.extend({
  pathForType() {
    return 'internal/counters/activity';
  },

  queryRecord(store, type, query) {
    let url = this.buildURL(type);
    // GET or POST ?
    // How to add namespace?
    return this.ajax(url, 'GET').then(response => {
      response.id = id;
      return response;
    });
  },
});
