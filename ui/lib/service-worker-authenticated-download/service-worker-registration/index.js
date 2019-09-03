import { addSuccessHandler, PROJECT_REVISION } from 'ember-service-worker/service-worker-registration';

function getToken() {
  // fix this later by allowing registration somewhere in the app lifecycle were we can have access to
  // services, etc.
  return Ember.Namespace.NAMESPACES_BY_ID['vault'].__container__.lookup('service:auth').currentToken;
}

addSuccessHandler(function(registration) {
  // attach the handler for the message event so we can send over the auth token
  navigator.serviceWorker.addEventListener('message', event => {
    let { action } = event.data;
    let port = event.ports[0];

    if (action === 'getToken') {
      let token = getToken();
      if (!token) {
        console.error('Unable to retrieve Vault tokent');
      }
      port.postMessage({ token: token });
    } else {
      console.error('Unknown event', event);
      port.postMessage({
        error: 'Unknown request',
      });
    }
  });

  // attempt to unregister the service worker on unload because we're not doing any sort of caching
  window.addEventListener('unload', function() {
    registration.unregister();
  });
});
