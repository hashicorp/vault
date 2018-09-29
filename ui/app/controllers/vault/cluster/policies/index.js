import { inject as service } from '@ember/service';
import { computed } from '@ember/object';
import Controller from '@ember/controller';

export default Controller.extend({
  flashMessages: service(),
  wizard: service(),

  queryParams: {
    page: 'page',
    pageFilter: 'pageFilter',
  },

  filter: null,
  page: 1,
  pageFilter: null,

  filterFocused: false,

  // set via the route `loading` action
  isLoading: false,

  filterMatchesKey: computed('filter', 'model', 'model.[]', function() {
    var filter = this.get('filter');
    var content = this.get('model');
    return !!(content && content.length && content.findBy('id', filter));
  }),

  firstPartialMatch: computed('filter', 'model', 'model.[]', 'filterMatchesKey', function() {
    var filter = this.get('filter');
    var content = this.get('model');
    if (!content) {
      return;
    }
    var filterMatchesKey = this.get('filterMatchesKey');
    var re = new RegExp('^' + filter);
    return filterMatchesKey
      ? null
      : content.find(function(key) {
          return re.test(key.get('id'));
        });
  }),

  actions: {
    setFilter: function(val) {
      this.set('filter', val);
    },
    setFilterFocus: function(bool) {
      this.set('filterFocused', bool);
    },
    deletePolicy(model) {
      let policyType = model.get('policyType');
      let name = model.id;
      let flash = this.get('flashMessages');
      model
        .destroyRecord()
        .then(() => {
          // this will clear the dataset cache on the store
          this.send('reload');
          flash.success(`${policyType.toUpperCase()} policy "${name}" was successfully deleted.`);
          if (this.get('wizard.featureState') === 'delete') {
            this.get('wizard').transitionFeatureMachine('delete', 'CONTINUE', policyType);
          }
        })
        .catch(e => {
          let errors = e.errors ? e.errors.join('') : e.message;
          flash.danger(
            `There was an error deleting the ${policyType.toUpperCase()} policy "${name}": ${errors}.`
          );
        });
    },
  },
});
