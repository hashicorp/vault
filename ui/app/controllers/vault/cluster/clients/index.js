import Controller from '@ember/controller';

export default class ClientsController extends Controller {
  queryParams = ['tab']; // ARG TODO remove
  tab = null;
}
