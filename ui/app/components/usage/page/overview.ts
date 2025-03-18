import { service } from '@ember/service';
import Component from '@glimmer/component';
import type UsageService from 'vault/services/usage';
export default class UsageOverviewPage extends Component {
  @service declare readonly usage: UsageService;
}
