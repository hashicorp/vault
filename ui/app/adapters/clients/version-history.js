import Application from '../application';

export default Application.extend({
  findAll() {
    return this.ajax(this.buildURL() + '/version-history', 'GET', {
      data: {
        list: true,
      },
    }).then((resp) => {
      return resp;
    });
  },
});
