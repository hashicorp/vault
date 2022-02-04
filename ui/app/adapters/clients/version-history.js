import Application from '../application';

export default Application.extend({
  queryRecord() {
    return this.ajax(this.buildURL() + '/version-history', 'GET', {
      data: {
        list: true,
      },
    }).then((resp) => {
      return resp;
    });
  },
});
