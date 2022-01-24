import Controller from '@ember/controller';

export default class ClientsController extends Controller {
  queryParams = ['tab', 'start', 'end']; // ARG TODO remove
  tab = null;
  start = null;
  end = null;
}
