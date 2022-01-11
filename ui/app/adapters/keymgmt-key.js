import ApplicationAdapter from './application';
import { encodePath } from 'vault/utils/path-encoding-helpers';

function pickKeys(obj, picklist) {
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

  _updateKey(backend, name, serialized) {
    let data = pickKeys(serialized, ['deletion_allowed', 'min_enabled_version']);
    return this.ajax(this.url(backend, name), 'PUT', { data });
  }

  _createKey(backend, name, serialized) {
    // Only type is allowed on create
    let data = pickKeys(serialized, ['type']);
    return this.ajax(this.url(backend, name), 'POST', { data });
  }

  async createRecord(store, type, snapshot) {
    const data = store.serializerFor(type.modelName).serialize(snapshot);
    const name = snapshot.attr('name');
    const backend = snapshot.attr('backend');
    // Keys must be created and then updated
    await this._createKey(backend, name, data);
    if (snapshot.attr('deletionAllowed')) {
      try {
        await this._updateKey(backend, name, data);
      } catch (e) {
        // TODO: Test how this works with UI
        throw new Error(`Key ${name} was created, but not all settings were saved`);
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
    try {
      const resp = await this.ajax(this.url(backend, name, 'PROVIDERS'), 'GET', {
        data: {
          list: true,
        },
      });
      // TODO: Get distribution
      return resp.data.keys ? resp.data.keys[0] : null;
    } catch (e) {
      if (e.httpStatus === 404) {
        return null;
      }
      // TODO: Handle no permissions
      console.error(e);
      return null;
    }
  }

  async queryRecord(store, type, query) {
    const { id, backend } = query;
    const keyData = await this.ajax(this.url(backend, id), 'GET');
    keyData.data.id = id;
    keyData.data.backend = backend;
    const provider = await this.getProvider(backend, id);
    return { ...keyData, provider };
  }

  async query(store, type, query) {
    const { backend } = query;
    return this.ajax(this.url(backend), 'GET', {
      data: {
        list: true,
      },
    }).then((res) => {
      res.backend = backend;
      return res;
    });
  }

  rotateKey(backend, id) {
    // TODO: re-fetch record data after
    return this.ajax(this.url(backend, id, 'ROTATE'), 'PUT');
  }
}
