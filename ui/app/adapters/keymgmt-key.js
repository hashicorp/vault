import ApplicationAdapter from './application';

export default class KeymgmtKeyAdapter extends ApplicationAdapter {
  namespace = 'v1';
  pathForType() {
    console.log('***** pathForType *******');
    return 'keymgmt/key';
  }

  query(store, type, query) {
    return super.query(store, type, query);
  }
}
