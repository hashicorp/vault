import SSHAdapter from './ssh';

export default SSHAdapter.extend({
  url(role) {
    return `/v1/${role.backend}/sign/${role.name}`;
  },
});
