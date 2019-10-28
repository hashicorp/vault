import Adapter from './pki';

export default Adapter.extend({
  url(role, snapshot) {
    if (snapshot.attr('signVerbatim') === true) {
      return `/v1/${role.backend}/sign-verbatim/${role.name}`;
    }
    return `/v1/${role.backend}/sign/${role.name}`;
  },

  pathForType() {
    return 'sign';
  },
});
