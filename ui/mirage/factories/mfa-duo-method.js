import { Factory } from 'ember-cli-mirage';

export default Factory.extend({
  api_hostname: 'api-foobar.duosecurity.com',
  mount_accessor: '',
  name: '', // returned but cannot be set at this time
  namespace_id: 'root',
  pushinfo: '',
  type: 'duo',
  use_passcode: false,
  username_template: '',

  afterCreate(record) {
    if (record.name) {
      console.warn('Endpoint ignored these unrecognized parameters: [name]'); // eslint-disable-line
      record.name = '';
    }
  },
});
