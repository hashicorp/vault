import { SELECTORS as ROLEFORM } from './pki-role-form';
import { SELECTORS as GENERATECERT } from './pki-role-generate';
import { SELECTORS as ISSUERDETAILS } from './pki-issuer-details';

export const SELECTORS = {
  breadcrumbContainer: '[data-test-breadcrumbs]',
  breadcrumbs: '[data-test-breadcrumbs] li',
  pageTitle: '[data-test-pki-role-page-title]',
  // TABS
  overviewTab: '[data-test-secret-list-tab="Overview"]',
  rolesTab: '[data-test-secret-list-tab="Roles"]',
  issuersTab: '[data-test-secret-list-tab="Issuers"]',
  certsTab: '[data-test-secret-list-tab="Certificates"]',
  keysTab: '[data-test-secret-list-tab="Keys"]',
  configTab: '[data-test-secret-list-tab="Configuration"]',
  // ROLES
  deleteRoleButton: '[data-test-pki-role-delete]',
  generateCertLink: '[data-test-pki-role-generate-cert]',
  signCertLink: '[data-test-pki-role-sign-cert]',
  editRoleLink: '[data-test-pki-role-edit-link]',
  createRoleLink: '[data-test-pki-role-create-link]',
  roleForm: {
    ...ROLEFORM,
  },
  generateCertForm: {
    ...GENERATECERT,
  },
  issuerDetails: {
    title: '[data-test-pki-issuer-page-title]',
    ...ISSUERDETAILS,
  },
};
