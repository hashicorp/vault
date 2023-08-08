import Component from '@glimmer/component';
import { inject as service } from '@ember/service';
import { inject as controller } from '@ember/controller';

export default class SidebarNavComponent extends Component {
  @service currentCluster;
  @service console;
  @controller('vault.cluster') clusterController;
}
