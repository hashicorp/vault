import { helper as buildHelper } from '@ember/component/helper';
import { capitalize } from '@ember/string';
import { assign } from '@ember/polyfills';

const SECRET_BACKENDS = {
  pki: {
    displayName: 'PKI',
    navigateTree: false,
    tabs: [
      {
        label: 'Overview',
        link: 'overview',
      },
    ],
  },
};

export function engineOptionsForBackend([backend, tab]) {
  const selected = SECRET_BACKENDS[backend];
  let backendOptions;

  if (selected && selected.tabs) {
    let tabData =
      selected.tabs.findBy('name', tab) || selected.tabs.findBy('modelPrefix', tab) || selected.tabs[0];
    backendOptions = assign({}, selected, tabData);
  } else if (selected) {
    backendOptions = selected;
  } else {
    backendOptions = assign({
      displayName: capitalize(backend),
    });
  }
  return backendOptions;
}

export default buildHelper(engineOptionsForBackend);
