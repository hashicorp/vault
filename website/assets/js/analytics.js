import { each } from './utils'

/* Segment's analytics.js provides a ready() function that is called once tracking is up and running */
/* Some clients block analytics.js, so to prevent errors, we assign noop functions if window.analytics isn't present */
window.analytics.ready(() => {
  const analytics = window.analytics || {
    trackLink: () => {},
    track: () => {},
    mock: true
  }

  // Track all button clicks
  track(
    '[data-ga-button]',
    el => {
      return {
        event: 'Click',
        category: 'Button',
        label: el.getAttribute('data-ga-button')
      }
    },
    true
  )

  // Track product subnav link clicks
  track(
    '[data-ga-product-subnav]',
    el => {
      return {
        event: 'Click',
        category: 'Product Subnav Navigation',
        label: el.getAttribute('data-ga-product-subnav')
      }
    },
    true
  )

  // Track meganav link clicks
  track(
    '[data-ga-meganav]',
    el => {
      return {
        event: 'Click',
        category: 'Meganav Navigation',
        label: el.getAttribute('data-ga-meganav')
      }
    },
    true
  )

  // Track footer link clicks
  track(
    '[data-ga-footer]',
    el => {
      return {
        event: 'Click',
        category: 'Footer Navigation',
        label: el.getAttribute('data-ga-footer')
      }
    },
    true
  )

  // Track outbound links
  track(
    'a[href^="http"]:not([href^="http://vaultproject.io"]):not([href^="https://vaultproject.io"]):not([href^="http://www.vaultproject.io"]):not([href^="https://www.vaultproject.io"])',
    el => {
      return {
        event: `Outbound Link | ${window.location.pathname}`,
        category: 'Outbound link',
        label: el.href
      }
    },
    true
  )

  // Note: Downloads are tracked from within the Product Downloader component

  /**
   * Wrapper for segment's track function that will track multiple elements,
   * normalize parameters, and easily switch between tracking links or events.
   * @param  {String} selector - query selector, multi element compatible
   * @param  {Function} cb - optional function that should return params, and will receive the element as a parameter
   * @param  {Boolean} [link=false] - if true, tracks a link click
   */
  function track(selector, cb, link = false) {
    each(document.querySelectorAll(selector), el => {
      let params = cb
      if (typeof cb === 'function') params = cb(el)
      const event = params.event
      delete params.event
      if (link) {
        analytics.trackLink(el, event, params)
      } else {
        el.addEventListener('click', () => {
          analytics.track(event, params)
        })
      }
    })
  }
})
