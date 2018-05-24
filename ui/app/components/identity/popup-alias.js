import Base from './_popup-base';

export default Base.extend({
  messageArgs(model) {
    let type = model.get('identityType');
    let id = model.id;
    return [type, id];
  },

  successMessage(type, id) {
    return `Successfully deleted ${type}: ${id}`;
  },

  errorMessage(e, type, id) {
    let error = e.errors ? e.errors.join(' ') : e.message;
    return `There was a problem deleting ${type}: ${id} - ${error}`;
  },

  transaction(model) {
    return model.destroyRecord();
  },
});
