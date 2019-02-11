const footer = require('@hashicorp/hashi-footer')
const nav = require('@hashicorp/hashi-nav')
const button = require('@hashicorp/hashi-button')
const megaNav = require('@hashicorp/hashi-mega-nav')
const productSubnav = require('@hashicorp/hashi-product-subnav')
const verticalTextBlockList = require('@hashicorp/hashi-vertical-text-block-list')
const sectionHeader = require('@hashicorp/hashi-section-header')
const content = require('@hashicorp/hashi-content')
const productDownloader = require('@hashicorp/hashi-product-downloader')
const docsSidenav = require('@hashicorp/hashi-docs-sidenav')
const hero = require('@hashicorp/hashi-hero')
const callouts = require('@hashicorp/hashi-callouts')
const splitCta = require('@hashicorp/hashi-split-cta')
const linkedTextSummaryList = require('@hashicorp/hashi-linked-text-summary-list')
const docsSitemap = require('@hashicorp/hashi-docs-sitemap')
const consentManager = require('@hashicorp/hashi-consent-manager')

module.exports = {
  'hashi-footer': footer,
  'hashi-nav': nav,
  'hashi-button': button,
  'hashi-docs-sidenav': docsSidenav,
  'hashi-mega-nav': megaNav,
  'hashi-product-subnav': productSubnav,
  'hashi-content': content,
  'hashi-product-downloader': productDownloader,
  'hashi-vertical-text-block-list': verticalTextBlockList,
  'hashi-section-header': sectionHeader,
  'hashi-hero': hero,
  'hashi-callouts': callouts,
  'hashi-split-cta': splitCta,
  'hashi-linked-text-summary-list': linkedTextSummaryList,
  'hashi-docs-sitemap': docsSitemap,
  'hashi-consent-manager': consentManager
}
