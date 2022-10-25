import Component from '@glimmer/component';
import { action } from '@ember/object';
import { next } from '@ember/runloop';

export default class SecretVersionMenu extends Component {
  onRefresh() {}

  @action
  closeDropdown(dropdown) {
    // strange issue where closing dropdown triggers full transition which redirects to auth screen in production builds
    // closing dropdown in next tick of run loop fixes it
    next(() => dropdown.actions.close());
  }
}
