import Component from '@glimmer/component';
import { action } from '@ember/object';

export default class DatabaseConnectionEdit extends Component {
  constructor(owner, args) {
    super(owner, args);
    const thisForm = this;

    // iterate over fields and set each one on the model for
    // form inputs
    args.model.allFields.forEach(field => {
      // iterate over all possible fields for the secret engine
      // and fill in either ember data value or default
      const matching = args.model[field];
      console.log(matching, `matches ${field}`);
      thisForm[field] = 'foo';
    });
  }

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
    console.log('submit', this.args.model);
  }

  @action
  updateValue(key, evt) {
    evt.preventDefault();
    console.log('update value', key, evt.target.value);
  }
}
