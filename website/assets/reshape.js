const button = require('@hashicorp/hashi-button')
const callouts = require('@hashicorp/hashi-callouts')
const caseStudySlider = require('@hashicorp/hashi-case-study-slider')
const consentManager = require('@hashicorp/hashi-consent-manager')
const content = require('@hashicorp/hashi-content')
const docsSidenav = require('@hashicorp/hashi-docs-sidenav')
const docsSitemap = require('@hashicorp/hashi-docs-sitemap')
const footer = require('@hashicorp/hashi-footer')
const hero = require('@hashicorp/hashi-hero')
const linkedTextSummaryList = require('@hashicorp/hashi-linked-text-summary-list')
const megaNav = require('@hashicorp/hashi-mega-nav')
const nav = require('@hashicorp/hashi-nav')
const productDownloader = require('@hashicorp/hashi-product-downloader')
const productSubnav = require('@hashicorp/hashi-product-subnav')
const sectionHeader = require('@hashicorp/hashi-section-header')
const splitCta = require('@hashicorp/hashi-split-cta')
const textAndContent = require('@hashicorp/hashi-text-and-content')
const verticalTextBlockList = require('@hashicorp/hashi-vertical-text-block-list')

const beforeAfterDiagram = require('./js/components/before-after-diagram')

module.exports = {
  'hashi-button': button,
  'hashi-callouts': callouts,
  'hashi-case-study-slider': caseStudySlider,
  'hashi-consent-manager': consentManager,
  'hashi-content': content,
  'hashi-docs-sidenav': docsSidenav,
  'hashi-docs-sitemap': docsSitemap,
  'hashi-footer': footer,
  'hashi-hero': hero,
  'hashi-linked-text-summary-list': linkedTextSummaryList,
  'hashi-mega-nav': megaNav,
  'hashi-nav': nav,
  'hashi-product-downloader': productDownloader,
  'hashi-product-subnav': productSubnav,
  'hashi-section-header': sectionHeader,
  'hashi-split-cta': splitCta,
  'hashi-text-and-content': textAndContent,
  'hashi-vertical-text-block-list': verticalTextBlockList,
  'hashi-before-after': beforeAfterDiagram
}
