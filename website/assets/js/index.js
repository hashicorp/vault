// components
import { each, initializeComponents } from './utils'
// external components
import nav from '@hashicorp/hashi-nav'
import footer from '@hashicorp/hashi-footer'
import newsletterSignupForm from '@hashicorp/hashi-newsletter-signup-form'
import docsSidebar from './components/docs-sidebar'

const components = initializeComponents({
  nav,
  footer,
  newsletterSignupForm,
  docsSidebar
})

const $ = document.querySelector.bind(document)
const $$ = document.querySelectorAll.bind(document)

// docs sidenav
const $sidebar = $('#sidebar')
if ($sidebar) {
  const $sidebarToggle = $sidebar.querySelector('.nav-toggle')
  const $sidebarUnderlay = $sidebar.querySelector('.underlay')

  $sidebarToggle.addEventListener('click', () => {
    $sidebar.classList.toggle('expanded')
  })
  $sidebarUnderlay.addEventListener('click', () => {
    $sidebar.classList.remove('expanded')
  })
}
