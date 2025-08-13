import { service } from '@ember/service';
import Component from '@glimmer/component';
import type Router from '@ember/routing/router';
import type FlashMessageService from 'vault/services/flash-messages';
import type ApiService from 'vault/services/api';
import SecretsEngineResource from 'vault/resources/secrets/engine';
interface Args {
  model: SecretsEngineResource;
}

export default class SecurityCard extends Component<Args> {
  @service declare readonly router: Router;
  @service declare readonly api: ApiService;
  @service declare readonly flashMessages: FlashMessageService;

  constructor(owner: unknown, args: Args) {
    super(owner, args);
  }
}
