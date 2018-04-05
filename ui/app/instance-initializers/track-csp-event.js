export function initialize(appInstance) {
  let service = appInstance.lookup('service:csp-event');
  service.attach();
}

export default {
  name: 'track-csp-event',
  initialize,
};
