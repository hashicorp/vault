import ApplicationAdapter from './application';

export default ApplicationAdapter.extend({
  query() {
    return this.ajax(this.urlForQuery(), 'GET');
  },

  urlForQuery() {
    return this.buildURL() + '/internal/ui/resultant-acl';
  },
});
