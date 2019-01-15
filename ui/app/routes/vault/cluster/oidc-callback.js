import { inject as service } from '@ember/service';
import Route from '@ember/routing/route';

export default Route.extend({
  store: service(),
  beforeModel() {
    let paramRegex = /\/ui(\/)?(.+)?\/vault\/oidc-callback\/(.+)$/;
    let params = window.location.pathname.match(paramRegex);
    // first will be the whole string, then the optional slash from the first capture, so we skip those
    let namespace = params[2];
    let path = params[3];
    let queryParams = window.location.search
      .substr(1)
      .split('&')
      .reduce(
        function(result, val) {
          let [keyName, keyVal] = val.split('=');
          result[keyName] = window.decodeURIComponent(keyVal);
          return result;
        },
        { namespace: namespace, path: path }
      );
    window.opener.postMessage(queryParams);
  },
});
