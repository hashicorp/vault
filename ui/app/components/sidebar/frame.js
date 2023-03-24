import Component from '@glimmer/component';
import { inject as service } from '@ember/service';

export default class SidebarNavComponent extends Component {
  @service currentCluster;
  @service console;
}
