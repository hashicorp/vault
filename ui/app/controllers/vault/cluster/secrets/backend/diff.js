import Controller from '@ember/controller';
import BackendCrumbMixin from 'vault/mixins/backend-crumb';
// import { inject as service } from '@ember/service';
import { action } from '@ember/object';
import { tracked } from '@glimmer/tracking';

export default class DiffController extends Controller.extend(BackendCrumbMixin) {
  // ARG TODO need to initialize the currentVersionData to currentVersion data, probably in the route and set on the controller
  @tracked
  currentVersionData = null;

  get getCompareVersion() {
    return this.model.currentVersion === 1 ? 0 : this.model.currentVersion - 1;
  }
  get currentVersion() {
    // another param that's selected this choose ssomething else from the action
    return this.model.currentVersion;
  }
  @action
  refreshModel() {
    this.send('refreshModel');
  }
  @action
  async selectCurrentVersion(selectedVersion) {
    let adapter = this.store.adapterFor('secret-v2-version');
    // return the version to json-edit and make it return the secret data
    // ARG TODO change secret path and make sure engineId changes from secret to other
    // ARG TODO options there's got to be a better way, i need to pass '["secret","my-secret","3"]'
    let string = `["${this.model.engineId}", "${this.model.id}", "${selectedVersion}"]`;
    let secretData = await adapter.querySecretDataByVersion(string);
    this.currentVersionData = secretData;
    // ARG TODO here stop here. secretData returns the data, send to  JSON edit consider later if new component.
  }
}
