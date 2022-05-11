import Controller from '@ember/controller';

export default class MfaEnforcementListController extends Controller {
  queryParams = {
    page: 'page',
  };

  page = 1;
}
