import Component from '@glimmer/component';

export default class NewFormFieldComponent extends Component {
  get value() {
    return this.args.model[this.args.name];
  }
}
