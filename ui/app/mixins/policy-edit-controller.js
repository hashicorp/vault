import { inject as service } from '@ember/service';
import Mixin from '@ember/object/mixin';

export default Mixin.create({
  flashMessages: service(),
  wizard: service(),
  actions: {
    deletePolicy(model) {
      const policyType = model.get('policyType');
      const name = model.get('name');
      const flash = this.flashMessages;
      model
        .destroyRecord()
        .then(() => {
          flash.success(`${policyType.toUpperCase()} policy "${name}" was successfully deleted.`);
          return this.transitionToRoute('vault.cluster.policies', policyType);
        })
        .catch((e) => {
          const errors = e.errors ? e.errors.join('') : e.message;
          flash.danger(
            `There was an error deleting the ${policyType.toUpperCase()} policy "${name}": ${errors}.`
          );
        });
    },

    savePolicy(model) {
      const flash = this.flashMessages;
      const policyType = model.get('policyType');
      const name = model.get('name');
      model
        .save()
        .then((m) => {
          flash.success(`${policyType.toUpperCase()} policy "${name}" was successfully saved.`);
          if (this.wizard.featureState === 'create') {
            this.wizard.transitionFeatureMachine('create', 'CONTINUE', policyType);
          }
          return this.transitionToRoute('vault.cluster.policy.show', m.get('policyType'), m.get('name'));
        })
        .catch(() => {
          // swallow error -- model.errors set by Adapter
          return;
        });
    },

    setModelName(model, e) {
      model.set('name', e.target.value.toLowerCase());
    },
  },
});
