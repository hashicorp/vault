import DS from 'ember-data';
import ApplicationSerializer from '../application';

export default ApplicationSerializer.extend(DS.EmbeddedRecordsMixin, {
  attrs: {
    aliases: { embedded: 'always' },
  },
});
