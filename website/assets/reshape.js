const footer = require('@hashicorp/hashi-footer')
const nav = require('@hashicorp/hashi-nav')
const button = require('@hashicorp/hashi-button')
const megaNav = require('@hashicorp/hashi-mega-nav')
const productSubnav = require('@hashicorp/hashi-product-subnav')
const verticalTextBlockList = require('@hashicorp/hashi-vertical-text-block-list')
const sectionHeader = require('@hashicorp/hashi-section-header')
const content = require('@hashicorp/hashi-content')
const productDownload = require('@hashicorp/hashi-product-downloader')
const docsSidebar = require('@hashicorp/hashi-docs-sidenav')

module.exports = {
  'hashi-footer': footer,
  'hashi-nav': nav,
  'hashi-button': button,
  'hashi-docs-sidebar': docsSidebar,
  'hashi-mega-nav': megaNav,
  'hashi-product-subnav': productSubnav,
  'hashi-content': content,
  'hashi-product-download': productDownload,
  'hashi-vertical-text-block-list': verticalTextBlockList,
  'hashi-section-header': sectionHeader
}
