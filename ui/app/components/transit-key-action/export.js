import Component from '@glimmer/component';
import { tracked } from '@glimmer/tracking';

export default class ExportComponent extends Component {
  @tracked
  wrapTTL = null;
  @tracked
  exportVersion = false;
}
