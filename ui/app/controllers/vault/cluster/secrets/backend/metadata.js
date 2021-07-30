import Controller from '@ember/controller';
import BackendCrumbMixin from 'vault/mixins/backend-crumb';
import { action } from '@ember/object';

export default class MetadataController extends Controller.extend(BackendCrumbMixin) {
  // ARG TODO not sure if I'm going to use this yet. Set the controller up for the mixin
  @action
  refreshModel() {
    this.send('refreshModel');
  }
}
