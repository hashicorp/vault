import { assign } from '@ember/polyfills';
import ApplicationAdapter from '../application';
import { allSettled } from 'rsvp';

export default ApplicationAdapter.extend({
  namespace: 'v1',
  pathForType() {
    // Confirm this isn't used by setting to bad val
    return 'jeepers';
  },

  urlFor(backend, id, type = 'dynamic') {
    let role = 'roles';
    if (type === 'static') {
      role = 'static-roles';
    }
    let url = `${this.buildURL()}/${backend}/${role}`;
    if (id) {
      url = `${this.buildURL()}/${backend}/${role}/${id}`;
    }
    return url;
  },

  staticRoles(backend, id) {
    return this.ajax(this.urlFor(backend, id, 'static'), 'GET', this.optionsForQuery(id));
  },

  dynamicRoles(backend, id) {
    return this.ajax(this.urlFor(backend, id), 'GET', this.optionsForQuery(id));
  },

  optionsForQuery(id) {
    let data = {};
    if (!id) {
      data['list'] = true;
    }
    return { data };
  },

  fetchByQuery(store, query) {
    const { backend, id } = query;
    return this.ajax(this.urlFor(backend, id), 'GET', this.optionsForQuery(id)).then(resp => {
      resp.id = id;
      resp.backend = backend;
      return resp;
    });
  },

  queryRecord(store, type, query) {
    const { backend, id } = query;
    const staticReq = this.staticRoles(backend, id);
    const dynamicReq = this.dynamicRoles(backend, id);

    return allSettled([staticReq, dynamicReq]).then(([staticResp, dynamicResp]) => {
      console.log(staticResp, dynamicResp);
      if (!staticResp.value && !dynamicResp.value) {
        // Throw error, both reqs failed
        throw dynamicResp.reason;
      }
      // Names are distinct across both types of role,
      // so only one request should ever come back with value
      let type = staticResp.value ? 'static' : 'dynamic';
      let successful = staticResp.value || dynamicResp.value;
      let resp = {
        data: {},
        backend,
        secret: id,
      };

      resp.data = assign({}, resp.data, successful.data, { backend, type, secret: id });

      return resp;
    });
  },

  query(store, type, query) {
    const { backend } = query;
    const staticReq = this.staticRoles(backend);
    const dynamicReq = this.dynamicRoles(backend);

    return allSettled([staticReq, dynamicReq]).then(([staticResp, dynamicResp]) => {
      // return hash(fetches).then((results) => {
      console.log('--------- FETCHES COMPLETE ---------');
      console.log(staticResp, dynamicResp);
      // if [].reason.httpStatus === 405, no permissions
      let resp = {
        backend,
        data: { keys: [] },
      };

      if (staticResp.reason && dynamicResp.reason) {
        console.log('Both failed');
        throw dynamicResp.reason;
      }
      // at least one request has data
      console.log('at least one successful req');
      let staticRoles = [];
      let dynamicRoles = [];

      if (staticResp.value) {
        staticRoles = staticResp.value.data.keys;
      }
      if (dynamicResp.value) {
        dynamicRoles = dynamicResp.value.data.keys;
      }

      resp.data = assign(
        {},
        resp.data,
        { keys: [...staticRoles, ...dynamicRoles] },
        { backend },
        { staticRoles, dynamicRoles }
      );

      console.log(resp, 'RESPONDING THIS');
      return resp;
    });
  },
});
