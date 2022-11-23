import Controller from '@ember/controller';
import { getOwner } from '@ember/application';
export default class PkiController extends Controller {
  get mountPoint() {
    return getOwner(this).mountPoint;
  }
}
