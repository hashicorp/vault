import Controller from '@ember/controller';
export default class MetricsController extends Controller {
  queryParams = ['tab', 'start', 'end'];
  tab = null;
  start = null;
  end = null;
}
