import Application from './application';

export default Application.extend({
  queryRecord() {
    return this.ajax(this.urlForQuery(), 'GET').then(resp => {
      resp.id = resp.request_id;
      let counters = resp.data.counters;
      counters.forEach(c => {
        const date = new Date(c.start_time);
        const month = date.getUTCMonth();
        const year = date.getUTCFullYear();
        // we have to manually create a new date with only the year and month
        // because the date is parsed in the users' local timezone which can
        // result in the month being off by 1
        const updated = new Date(year, month);
        c.start_time = updated;
      });
      return resp;
    });
  },

  urlForQuery() {
    return this.buildURL() + '/internal/counters/requests';
  },
});
