import Controller from '@ember/controller';

export default class MfaMethodsListController extends Controller {
  queryParams = {
    page: 'page',
  };

  page = 1;
}
