import Component from '@glimmer/component';
import errorMessage from 'vault/utils/error-message';

export default class NewFormComponent extends Component {
  get errorMessage() {
    console.log(this.args.formError);
    return errorMessage(this.args.formError);
  }
}
