import ApplicationSerializer from './application';

export default class SecretListSerializer extends ApplicationSerializer {
  normalizeResponse(store, primaryModelClass, payload, id, requestType) {
    // queryRecord will already have set this, and we won't have an id here
    if (payload.data.keys) {
      payload.data.secrets = payload.data.keys.map((key) => ({ id: key }));
      delete payload.data.keys;
      payload.id = id || `${payload.backend}-list`;
      // return [
      //   {
      //     id: 'record-id',
      //     data: {
      //       secrets: payload.data.keys,
      //     },
      //   },
      // ];
    }
    return super.normalizeResponse(store, primaryModelClass, payload, id, requestType);
  }
}
