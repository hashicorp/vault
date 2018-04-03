import Ember from 'ember';

let { inject } = Ember;

export default Ember.Mixin.create({
  flashMessages: inject.service(),
  actions: {
    deletePolicy(model) {
      let policyType = model.get('policyType');
      let name = model.get('name');
      let flash = this.get('flashMessages');
      model
        .destroyRecord()
        .then(() => {
          flash.success(`${policyType.toUpperCase()} policy "${name}" was successfully deleted.`);
          return this.transitionToRoute('vault.cluster.policies', policyType);
        })
        .catch(e => {
          let errors = e.errors ? e.errors.join('') : e.message;
          flash.danger(
            `There was an error deleting the ${policyType.toUpperCase()} policy "${name}": ${errors}.`
          );
        });
    },

    savePolicy(model) {
      let flash = this.get('flashMessages');
      let policyType = model.get('policyType');
      let name = model.get('name');
      model.save().then(m => {
        flash.success(`${policyType.toUpperCase()} policy "${name}" was successfully saved.`);
        return this.transitionToRoute('vault.cluster.policy.show', m.get('policyType'), m.get('name'));
      });
    },

    setModelName(model, e) {
      model.set('name', e.target.value.toLowerCase());
    },
  },
});
