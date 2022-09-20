import Controller from '@ember/controller';
import { action } from '@ember/object';

export default class PkiCertificatesIndexController extends Controller {
  @action setFilter(val) {
    this.filter = val;
  }
  @action setFilterFocus(bool) {
    this.filterFocused = bool;
  }
}
