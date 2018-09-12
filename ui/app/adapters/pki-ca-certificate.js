import ApplicationAdapter from './application';

export default ApplicationAdapter.extend({
  namespace: 'v1',

  url(snapshot, action) {
    const { backend, caType, type } = snapshot.attributes();
    if (action === 'sign-intermediate') {
      return `/v1/${backend}/root/sign-intermediate`;
    }
    if (action === 'set-signed-intermediate') {
      return `/v1/${backend}/intermediate/set-signed`;
    }
    if (action === 'upload') {
      return `/v1/${backend}/config/ca`;
    }
    return `/v1/${backend}/${caType}/generate/${type}`;
  },

  createRecordOrUpdate(store, type, snapshot, requestType) {
    const serializer = store.serializerFor('application');
    const isUpload = snapshot.attr('uploadPemBundle');
    const isSetSignedIntermediate = snapshot.adapterOptions.method === 'setSignedIntermediate';
    let action = snapshot.adapterOptions.method === 'signIntermediate' ? 'sign-intermediate' : null;
    let data;
    if (isUpload) {
      action = 'upload';
      data = { pem_bundle: snapshot.attr('pemBundle') };
    } else if (isSetSignedIntermediate) {
      action = 'set-signed-intermediate';
      data = { certificate: snapshot.attr('certificate') };
    } else {
      data = serializer.serialize(snapshot, requestType);
    }

    return this.ajax(this.url(snapshot, action), 'POST', { data }).then(response => {
      // uploading CA, setting signed intermediate cert, and attempting to generate
      // a new CA if one exists, all return a 204
      if (!response) {
        response = {};
      }
      response.id = snapshot.id;
      response.modelName = type.modelName;
      store.pushPayload(type.modelName, response);
    });
  },

  createRecord() {
    return this.createRecordOrUpdate(...arguments);
  },

  updateRecord() {
    return this.createRecordOrUpdate(...arguments);
  },

  deleteRecord(store, type, snapshot) {
    const backend = snapshot.attr('backend');
    return this.ajax(`/v1/${backend}/root`, 'DELETE');
  },
});
