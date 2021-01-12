import { helper as buildHelper } from '@ember/component/helper';
import { capitalize } from '@ember/string';
import { assign } from '@ember/polyfills';

const DEFAULT_DISPLAY = {
  searchPlaceholder: 'Filter secrets',
  item: 'secret',
  create: 'Create secret',
  navigateTree: true,
  editComponent: 'secret-edit',
  listItemPartial: 'partials/secret-list/item',
};
const SECRET_BACKENDS = {
  aws: {
    displayName: 'AWS',
    searchPlaceholder: 'Filter roles',
    item: 'role',
    create: 'Create role',
    navigateTree: false,
    editComponent: 'role-aws-edit',
    listItemPartial: 'partials/secret-list/aws-role-item',
  },
  pki: {
    displayName: 'PKI',
    navigateTree: false,
    listItemPartial: 'partials/secret-list/pki-role-item',
    tabs: [
      {
        name: 'roles',
        label: 'Roles',
        searchPlaceholder: 'Filter roles',
        item: 'role',
        create: 'Create role',
        editComponent: 'role-pki-edit',
      },
      {
        name: 'certs',
        modelPrefix: 'cert/',
        label: 'Certificates',
        searchPlaceholder: 'Filter certificates',
        item: 'certificates',
        create: 'Create role',
        tab: 'certs',
        listItemPartial: 'partials/secret-list/pki-cert-item',
        editComponent: 'pki-cert-show',
      },
    ],
  },
  ssh: {
    displayName: 'SSH',
    searchPlaceholder: 'Filter roles',
    item: 'role',
    create: 'Create role',
    navigateTree: false,
    editComponent: 'role-ssh-edit',
    listItemPartial: 'partials/secret-list/ssh-role-item',
  },
  database: {
    displayName: 'Database',
    navigateTree: false,
    listItemPartial: 'partials/secret-list/database-list-item',
    tabs: [
      // we have the landing route Overview, and tabs connection and roles
      {
        name: 'overview',
        modelPrefix: 'overview/',
        label: 'Overview',
        item: 'overview',
        tab: 'overview',
        editComponent: 'database-overview-edit',
      },
      {
        name: 'connection',
        modelPrefix: 'connection/',
        label: 'Connections',
        searchPlaceholder: 'Filter connections',
        item: 'connection',
        create: 'Create connection',
        tab: 'connection',
        editComponent: 'database-connection-edit',
      },
      {
        name: 'role',
        modelPrefix: 'role/',
        label: 'Roles',
        searchPlaceholder: 'Filter roles',
        item: 'role',
        create: 'Create role',
        tab: 'role',
        editComponent: 'database-role-edit',
      },
    ],
  },
  transform: {
    displayName: 'Transformation',
    navigateTree: false,
    listItemPartial: 'partials/secret-list/transform-list-item',
    tabs: [
      {
        name: 'transformations',
        label: 'Transformations',
        searchPlaceholder: 'Filter transformations',
        item: 'transformation',
        create: 'Create transformation',
        editComponent: 'transformation-edit',
        listItemPartial: 'partials/secret-list/transform-transformation-item',
      },
      {
        name: 'role',
        modelPrefix: 'role/',
        label: 'Roles',
        searchPlaceholder: 'Filter roles',
        item: 'role',
        create: 'Create role',
        tab: 'role',
        editComponent: 'transform-role-edit',
      },
      {
        name: 'template',
        modelPrefix: 'template/',
        label: 'Templates',
        searchPlaceholder: 'Filter templates',
        item: 'template',
        create: 'Create template',
        tab: 'template',
        editComponent: 'transform-template-edit',
      },
      {
        name: 'alphabet',
        modelPrefix: 'alphabet/',
        label: 'Alphabets',
        searchPlaceholder: 'Filter alphabets',
        item: 'alphabet',
        create: 'Create alphabet',
        tab: 'alphabet',
        editComponent: 'alphabet-edit',
      },
    ],
  },
  transit: {
    searchPlaceholder: 'Filter keys',
    item: 'key',
    create: 'Create encryption key',
    navigateTree: false,
    editComponent: 'transit-edit',
    listItemPartial: 'partials/secret-list/item',
  },
};

export function optionsForBackend([backend, tab]) {
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
      displayName: backend === 'kv' ? 'KV' : capitalize(backend),
    });
  }
  return backendOptions;
}

export default buildHelper(optionsForBackend);
