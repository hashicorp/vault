import ApplicationAdapter from './application';
import DS from 'ember-data';
import Ember from 'ember';

export default ApplicationAdapter.extend({
  namespace: 'v1',

  defaultSerializer: 'config',

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
    const serializer = store.serializerFor(this.get('defaultSerializer'));
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
    return this.ajax(url, 'POST', { data });
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
      Ember.set(error, 'httpStatus', 404);
      throw error;
    }
    return this[`fetch${Ember.String.capitalize(section)}`](backendPath);
  },

  id(backendPath) {
    return backendPath + '-config-ca';
  },

  fetchCert(backendPath) {
    // these are all un-authed so using `fetch` directly works
    const derURL = `/v1/${backendPath}/ca`;
    const pemURL = `${derURL}/pem`;
    const chainURL = `${derURL}_chain`;

    return Ember.RSVP.hash({
      backend: backendPath,
      id: this.id(backendPath),
      der: this.rawRequest(derURL, { unauthenticate: true }).then(response => response.blob()),
      pem: this.rawRequest(pemURL, { unauthenticate: true }).then(response => response.text()),
      ca_chain: this.rawRequest(chainURL, { unauthenticate: true }).then(response => response.text()),
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
          return Ember.RSVP.resolve({ id });
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
    return Ember.RSVP.resolve({ id, backend: backendPath });
  },

  queryRecord(store, type, query) {
    const { backend, section } = query;
    return this.fetchSection(backend, section).then(resp => {
      resp.backend = backend;
      return resp;
    });
  },
});
