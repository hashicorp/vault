import Base from './_popup-base';

export default Base.extend({
  messageArgs(model) {
    const type = model.get('identityType');
    const id = model.id;
    return [type, id];
  },

  successMessage(type, id) {
    return `Successfully deleted ${type}: ${id}`;
  },

  errorMessage(e, type, id) {
    const error = e.errors ? e.errors.join(' ') : e.message;
    return `There was a problem deleting ${type}: ${id} - ${error}`;
  },

  transaction(model) {
    return model.destroyRecord();
  },
});
