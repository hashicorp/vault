export function initialize(appInstance) {
  let service = appInstance.lookup('service:permissions');
  service.checkAuthToken.perform();
}

export default {
  name: 'start-permissions-task',
  initialize,
};
