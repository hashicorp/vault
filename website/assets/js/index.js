// components
import { each, initializeComponents } from './utils'
// external components
import nav from '@hashicorp/hashi-nav'
import footer from '@hashicorp/hashi-footer'
import newsletterSignupForm from '@hashicorp/hashi-newsletter-signup-form'
import productSubnav from '@hashicorp/hashi-product-subnav'
import megaNav from '@hashicorp/hashi-mega-nav'
import productDownloader from '@hashicorp/hashi-product-downloader'
import hero from '@hashicorp/hashi-hero'
import docsSidenav from '@hashicorp/hashi-docs-sidenav'
import consentManager from '@hashicorp/hashi-consent-manager'

const components = initializeComponents({
  nav,
  footer,
  newsletterSignupForm,
  productSubnav,
  megaNav,
  productDownloader,
  hero,
  docsSidenav,
  consentManager
})
