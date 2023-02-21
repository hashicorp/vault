import Controller from '@ember/controller';
import { getOwner } from '@ember/application';

export default class PkiRolesIndexController extends Controller {
  get mountPoint() {
    return getOwner(this).mountPoint;
  }
}
