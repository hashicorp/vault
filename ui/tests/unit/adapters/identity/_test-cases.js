export const storeMVP = {
  serializerFor() {
    return {
      serializeIntoHash() {},
    };
  },
};

export default function(modelName) {
  return [
    {
      adapterMethod: 'findRecord',
      args: [null, { modelName }, 'foo'],
      url: `/v1/${modelName}/id/foo`,
      method: 'GET',
    },

    {
      adapterMethod: 'createRecord',
      args: [storeMVP, { modelName }],
      url: `/v1/${modelName}`,
      method: 'POST',
    },
    {
      adapterMethod: 'updateRecord',
      args: [storeMVP, { modelName }, { id: 'foo' }],
      url: `/v1/${modelName}/id/foo`,
      method: 'PUT',
    },
    {
      adapterMethod: 'deleteRecord',
      args: [storeMVP, { modelName }, { id: 'foo' }],
      url: `/v1/${modelName}/id/foo`,
      method: 'DELETE',
    },
    {
      adapterMethod: 'query',
      args: [null, { modelName }, {}],
      url: `/v1/${modelName}/id?list=true`,
      method: 'GET',
    },
  ];
}
