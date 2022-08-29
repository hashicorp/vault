import { helper as buildHelper } from '@ember/component/helper';
import { capitalize } from '@ember/string';
import { assign } from '@ember/polyfills';

// ARG TODO confirm what you need on default display
const DEFAULT_DISPLAY = {
  searchPlaceholder: 'Filter secrets',
  item: 'secret',
  create: 'Create secret',
  navigateTree: true,
  editComponent: 'secret-edit',
  listItemPartial: 'secret-list/item',
};
const SECRET_BACKENDS = {
  pki: {
    displayName: 'PKI',
    navigateTree: false,
    tabs: [
      {
        label: 'Overview',
        link: 'overview',
      },
      // {
      //   name: 'cert',
      //   modelPrefix: 'cert/',
      //   label: 'Certificates',
      //   searchPlaceholder: 'Filter certificates',
      //   item: 'certificates',
      //   create: 'Create role',
      //   tab: 'cert',
      //   listItemPartial: 'secret-list/pki-cert-item',
      //   editComponent: 'pki/pki-cert-show',
      // },
    ],
  },
};

export function optionsForBackend([backend, tab]) {
  // ARG TODO sort through this and see what functionality you can clean up specific for PKI only
  const selected = SECRET_BACKENDS[backend];
  let backendOptions;

  if (selected && selected.tabs) {
    let tabData =
      selected.tabs.findBy('name', tab) || selected.tabs.findBy('modelPrefix', tab) || selected.tabs[0];
    backendOptions = assign({}, selected, tabData);
  } else if (selected) {
    backendOptions = selected;
  } else {
    backendOptions = assign({}, DEFAULT_DISPLAY, {
      displayName: capitalize(backend),
    });
  }
  return backendOptions;
}

export default buildHelper(optionsForBackend);
