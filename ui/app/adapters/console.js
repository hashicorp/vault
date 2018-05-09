import ApplicationAdapter from './application';

export default ApplicationAdapter.extend({
  namespace: 'v1',
  pathForType(modelName) {
    return modelName;
  },
});
