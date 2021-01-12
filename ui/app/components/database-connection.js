import Component from '@glimmer/component';
import { action } from '@ember/object';

// Template grabs model from ember-data

export default class DatabaseConnectionEdit extends Component {
  constructor(owner, args) {
    super(owner, args);
    const thisForm = this;
    console.log('<<<< CONSTRUCTED ARGS', args.model.fieldAttrs);
    args.model.allFields.forEach(field => {
      // iterate over all possible fields for the secret engine
      // and fill in either ember data value or default
      const matching = args.model[field];
      console.log(matching, `matches ${field}`);
      thisForm[field] = 'foo';
    });
  }

  key = '';
  // mode = 'show';
  // tab
  // model
  // mode
  // root
  // capabilities
  // onRefresh
  // onToggleAdvancedEdit
  // initialKey
  // baseKey
  // preferAdvancedEdit
  @action
  async handleSubmit(evt) {
    evt.preventDefault();
    // this.args.sendMessage(this.body);
    // this.body = '';
    console.log('submit');
  }

  @action
  updateValue(key, evt) {
    evt.preventDefault();
    console.log('update value', key, evt.target.value);
  }
}
