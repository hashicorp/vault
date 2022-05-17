import Controller from '@ember/controller';
import { inject as service } from '@ember/service';
import { tracked } from '@glimmer/tracking';
import { action } from '@ember/object';
import { capitalize } from '@ember/string';

export default class MfaMethodCreateController extends Controller {
  @service store;

  queryParams = ['type'];
  methodNames = ['TOTP', 'Duo', 'Okta', 'PingID'];

  @tracked type = null;
  @tracked selectedType;
  @tracked method;

  get formattedSelectedType() {
    if (!this.selectedType) return '';
    return this.selectedType === 'totp' ? this.selectedType.toUpperCase() : capitalize(this.selectedType);
  }

  @action
  createMethod() {
    this.method = this.store.createRecord('mfa-method', { type: this.selectedType });
    this.type = this.selectedType; // set selectedType to query param for state tracking
  }
}
