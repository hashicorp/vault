import { hash, resolve } from 'rsvp';
import { capitalize } from '@ember/string';
import { set } from '@ember/object';
import ApplicationAdapter from './application';
import DS from 'ember-data';

export default ApplicationAdapter.extend({
  namespace: 'v1',

  urlFor(backend, section) {
    const urls = {
      tidy: `/v1/${backend}/tidy`,
      urls: `/v1/${backend}/config/urls`,
      crl: `/v1/${backend}/config/crl`,
    };
    return urls[section];
  },

  createOrUpdate(store, type, snapshot) {
    const url = this.urlFor(snapshot.record.get('backend'), snapshot.adapterOptions.method);
    const serializer = store.serializerFor(type.modelName);
    if (!url) {
      return;
    }
    const data = snapshot.adapterOptions.fields.reduce((data, field) => {
      let attr = snapshot.attr(field);
      if (attr) {
        serializer.serializeAttribute(snapshot, data, field, attr);
      } else {
        data[serializer.keyForAttribute(field)] = attr;
      }
      return data;
    }, {});
    return this.ajax(url, 'POST', { data }).then(resp => {
      let response = resp || {};
      response.id = `${snapshot.record.get('backend')}-${snapshot.adapterOptions.method}`;
      return response;
    });
  },

  createRecord() {
    return this.createOrUpdate(...arguments);
  },

  updateRecord() {
    return this.createOrUpdate(...arguments, 'update');
  },

  fetchSection(backendPath, section) {
    const sections = ['cert', 'urls', 'crl', 'tidy'];
    if (!section || !sections.includes(section)) {
      const error = new DS.AdapterError();
      set(error, 'httpStatus', 404);
      throw error;
    }
    return this[`fetch${capitalize(section)}`](backendPath);
  },

  id(backendPath) {
    return backendPath + '-config-ca';
  },

  fetchCert(backendPath) {
    // these are all un-authed so using `fetch` directly works
    const derURL = `/v1/${backendPath}/ca`;
    const pemURL = `${derURL}/pem`;
    const chainURL = `${derURL}_chain`;

    return hash({
      backend: backendPath,
      id: this.id(backendPath),
      der: this.rawRequest(derURL, 'GET', { unauthenticated: true }).then(response => response.blob()),
      pem: this.rawRequest(pemURL, 'GET', { unauthenticated: true }).then(response => response.text()),
      ca_chain: this.rawRequest(chainURL, 'GET', { unauthenticated: true }).then(response => response.text()),
    });
  },

  fetchUrls(backendPath) {
    const url = `/v1/${backendPath}/config/urls`;
    const id = this.id(backendPath);
    return this.ajax(url, 'GET')
      .then(resp => {
        resp.id = id;
        resp.backend = backendPath;
        return resp;
      })
      .catch(e => {
        if (e.httpStatus === 404) {
          return resolve({ id });
        } else {
          throw e;
        }
      });
  },

  fetchCrl(backendPath) {
    const url = `/v1/${backendPath}/config/crl`;
    const id = this.id(backendPath);
    return this.ajax(url, 'GET')
      .then(resp => {
        resp.id = id;
        resp.backend = backendPath;
        return resp;
      })
      .catch(e => {
        if (e.httpStatus === 404) {
          return { id };
        } else {
          throw e;
        }
      });
  },

  fetchTidy(backendPath) {
    const id = this.id(backendPath);
    return resolve({ id, backend: backendPath });
  },

  queryRecord(store, type, query) {
    const { backend, section } = query;
    return this.fetchSection(backend, section).then(resp => {
      resp.backend = backend;
      return resp;
    });
  },
});
