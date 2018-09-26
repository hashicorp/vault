// components
import { each, initializeComponents } from './utils'
// external components
import nav from '@hashicorp/hashi-nav'
import footer from '@hashicorp/hashi-footer'
import newsletterSignupForm from '@hashicorp/hashi-newsletter-signup-form'
import productSubnav from '@hashicorp/hashi-product-subnav'
import docsSidebar from './components/docs-sidebar'
import megaNav from '@hashicorp/hashi-mega-nav'

const components = initializeComponents({
  nav,
  footer,
  newsletterSignupForm,
  docsSidebar,
  productSubnav,
  megaNav
})
