import Component from '@glimmer/component';
import { action } from '@ember/object';
import { inject as service } from '@ember/service';

export default class SecretVersionMenu extends Component {
  @service router;

  onRefresh() {}

  @action
  changeVersion(version, dropdown) {
    dropdown.actions.close();
    this.router.transitionTo({
      queryParams: { version },
    });
  }
}
