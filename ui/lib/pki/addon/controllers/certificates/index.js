import PkiController from '../pki';
import { action } from '@ember/object';

export default class PkiCertificatesIndexController extends PkiController {
  @action setFilter(val) {
    this.filter = val;
  }
  @action setFilterFocus(bool) {
    this.filterFocused = bool;
  }
}
