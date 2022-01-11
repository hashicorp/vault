import ApplicationAdapter from './application';
import { encodePath } from 'vault/utils/path-encoding-helpers';

function pickKeys(obj, picklist) {
  console.log('Picking keys from', obj, picklist);
  const data = {};
  return Object.keys(obj).forEach((key) => {
    if (picklist.indexOf(key) >= 0) {
      data[key] = obj[key];
    }
  });
}
export default class KeymgmtKeyAdapter extends ApplicationAdapter {
  namespace = 'v1';
  pathForType() {
    console.log('***** pathForType *******');
    return 'keymgmt/key';
  }

  url(backend, id, type) {
    const url = `${this.buildURL()}/${backend}/key`;
    if (id) {
      if (type === 'ROTATE') {
        return url + '/' + encodePath(id) + '/rotate';
      } else if (type === 'PROVIDERS') {
        return url + '/' + encodePath(id) + '/kms';
      }
      return url + '/' + encodePath(id);
    }
    return url;
  }

  updateKey(backend, name, serialized) {
    let data = pickKeys(serialized, ['deletion_allowed', 'min_enabled_version']);
    return this.ajax(this.url(backend, name), 'PUT', { data });
  }

  createKey(backend, name, serialized) {
    // Only type is allowed on create
    let data = pickKeys(serialized, ['type']);
    return this.ajax(this.url(backend, name), 'POST', { data });
  }

  async createRecord(store, type, snapshot) {
    const data = store.serializerFor(type.modelName).serialize(snapshot);
    const name = snapshot.attr('name');
    const backend = snapshot.attr('backend');
    // Keys must be created and then updated
    await this.createKey(backend, name, data);
    if (snapshot.attr('deletionAllowed')) {
      try {
        await this.updateKey(backend, name, data);
      } catch (e) {
        console.debug(e);
        throw new Error(`Key ${name} was created, but not all settings were saved. ${e.message}`);
      }
    }
    return {
      data: {
        ...data,
        id: name,
        backend,
      },
    };
  }

  async getProvider(backend, name) {
    // TODO: Handle no permissions
    const resp = await this.ajax(this.url(backend, name, 'PROVIDERS'), 'GET');
    // TODO: Get distribution
    return resp.data.keys ? resp.data.keys[0] : null;
  }

  async queryRecord(store, type, query) {
    const { id, backend } = query;
    const keyData = await this.ajax(this.url(backend, id), 'GET');
    keyData.data.id = id;
    keyData.data.backend = backend;
    const provider = await this.getProvider(backend, id);
    return { ...keyData, provider };
  }
}
