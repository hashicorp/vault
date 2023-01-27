import ApplicationSerializer from '../application';

export default class PkiRoleSerializer extends ApplicationSerializer {
  attrs = {
    name: { serialize: false },
  };
}
