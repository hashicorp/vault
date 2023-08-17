import { addSuccessHandler } from 'ember-service-worker/service-worker-registration';

addSuccessHandler(function (registration) {
  // attempt to unregister the service worker on unload because we're not doing any sort of caching
  window.addEventListener('unload', function () {
    registration.unregister();
  });
});
