import ApplicationAdapter from '../application';
import { all } from 'rsvp';

export default class KeymgmtKeyAdapter extends ApplicationAdapter {
  namespace = 'v1';
  listPayload = { data: { list: true } };

  pathForType() {
    return 'keymgmt/kms';
  }
  async createRecord(store, type, snapshot) {
    // create uses PUT instead of POST
    const data = store.serializerFor(type.modelName).serialize(snapshot);
    const url = `${this.buildURL(type.modelName)}/${snapshot.attr('name')}`;
    return this.ajax(url, 'PUT', { data }).then(() => data);
  }
  findRecord(store, type, name) {
    return super.findRecord(...arguments).then((resp) => {
      resp.data = { ...resp.data, name };
      return resp;
    });
  }
  async query(store, type) {
    return this.ajax(this.buildURL(type.modelName), 'GET', this.listPayload).then(async (resp) => {
      // additional data is needed to fullfil the list view requirements
      // pull in full record for listed items
      const records = await all(resp.data.keys.map((name) => this.findRecord(store, type, name)));
      resp.data.keys = records.map((record) => record.data);
      return resp;
    });
  }
  async queryRecord(store, type, query) {
    return this.findRecord(store, type, query.id);
  }
}
