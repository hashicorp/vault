import { render } from 'preact'
import { hydrateInitialState } from 'reshape-preact-components/lib/browser'

// rehydrates and initializes top-level preact components
export function initializeComponents(obj) {
  const res = {}

  for (let k in obj) {
    const name = getName(k)
    res[name] = []
    each(document.querySelectorAll(`.g-${name}`), el => {
      // do not initialize nested components
      const matches = Object.keys(obj)
        .map(getName)
        .reduce((m, name) => {
          const parent = findParent(el, `.g-${name}`)
          if (parent) m.push(parent)
          return m
        }, [])
      if (matches.length > 1) return
      // if there's no data-state, don't try
      if (!el.dataset.state || !el.dataset.state.length) {
        return
      }
      // otherwise, initialize away
      const vdom = hydrateInitialState(el.dataset.state, {
        [`hashi-${name}`]: obj[k]
      })

      res[name].push(render(vdom, el.parentElement, el))
    })
  }

  return res

  function getName(s) {
    return s.replace(/([A-Z])/g, '-$1').toLowerCase()
  }
}

// iterates through a NodeList
export function each(list, cb) {
  for (let i = 0; i < list.length; i++) {
    cb(list[i], i)
  }
}

// polyfills object-fit in unsupported browsers
export function fixObjectFit() {
  if (Modernizr.objectfit) {
    import('object-fit-images').then(ofi => {
      ofi.default()
    })
  }
}

// given an element and selector, finds the closest parent element. doesn't
// handle attribute selectors, just class, id, and element name
export function findParent(el, selector) {
  const firstChar = selector[0]
  if (firstChar === '.') {
    if (el.classList.contains(selector.substr(1))) return el
  } else if (firstChar === '#') {
    if (el.id === selector.substr(1)) return el
  } else {
    if (el.tagName.toLowerCase() === selector) return el
  }
  if (!el.parentNode.tagName) return undefined
  return findParent(el.parentNode, selector)
}
