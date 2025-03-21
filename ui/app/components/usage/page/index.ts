import Component from '@ember/component';
import { service } from '@ember/service';
import UsageService from 'vault/services/usage';

export default class ClientsActivityComponent extends Component {
  @service declare readonly usage: UsageService;
}
