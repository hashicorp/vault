import Application from './application';

export default Application.extend({
  queryRecord() {
    return this.ajax(this.urlForQuery(), 'GET').then(resp => {
      resp.id = resp.request_id;
      // resp.data.totalTokens = resp.data.counters.service_tokens.total;
      return resp;
    });
  },

  urlForQuery() {
    return this.buildURL() + '/internal/counters/tokens';
  },
});
