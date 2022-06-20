import Controller from '@ember/controller';
import { inject as service } from '@ember/service';
import { tracked } from '@glimmer/tracking';
import { action } from '@ember/object';

export default class OidcClientCreateController extends Controller {
  @service store;

  @tracked enforcement;
  @tracked enforcementPreference = 'new';
  @tracked radioSelectedValue;

  @action
  onAssignAccessChange(preference) {
    this.radioSelectedValue = preference;
    // TODO will need to send this value to model on submit.
  }
}
