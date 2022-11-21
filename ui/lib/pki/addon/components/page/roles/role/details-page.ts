import { action } from '@ember/object';
import Component from '@glimmer/component';

// TODO: pull this in from route model once it's TS
interface Args {
  role: {
    backend: string;
    id: string;
  };
}

export default class DetailsPage extends Component<Args> {
  get breadcrumbs() {
    return [
      { label: this.args.role.backend || 'pki', path: 'overview' },
      { label: 'roles', path: 'roles.index' },
      { label: this.args.role.id },
    ];
  }

  @action deleteRole() {
    // TODO: delete role
  }
}
