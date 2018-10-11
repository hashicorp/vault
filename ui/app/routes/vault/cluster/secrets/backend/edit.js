import EditBase from './secret-edit';

export default EditBase.extend({
  queryParams: {
    version: {
      refreshModel: true,
    },
  },
});
