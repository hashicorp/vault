import Application from './application';

export default Application.extend({
  queryRecord() {
    return this.ajax(this.urlForQuery(), 'GET').then(resp => {
      resp.id = resp.request_id;
      return resp;
    });
  },

  urlForQuery() {
    return this.buildURL() + '/internal/counters/requests';
  },
});
