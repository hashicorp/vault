import Controller from '@ember/controller';

export default class KvSecretDetailsController extends Controller {
  queryParams = ['version'];
  version = null;
}
