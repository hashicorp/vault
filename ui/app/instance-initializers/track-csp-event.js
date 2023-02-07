export function initialize(appInstance) {
  const service = appInstance.lookup('service:csp-event');
  service.attach();
}

export default {
  name: 'track-csp-event',
  initialize,
};
