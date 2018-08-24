const { h, Component } = require('preact')
const { decode } = require('reshape-preact-components')
const assign = require('object-assign')

module.exports = class Sidebar extends Component {
  render() {
    const current = this.props.current_page.split('/').slice(1)
    const category = this.props.category
    const order = decode(this.props.order)
    const data = decode(this.props.data).map(p => {
      p.path = p.path
        .split('/')
        .slice(1)
        .join('/')
      return p
    })

    return (
      <div data-state={this.props._state} class="g-docs-sidebar">
        <ul class="nav docs-nav">
          {this.renderNavTree(
            category,
            this.matchOrderToPageData(order, data),
            current
          )}
        </ul>
      </div>
    )
  }

  // replace all terminal page nodes with page data from middleman
  matchOrderToPageData(content, pageData) {
    // go through each item in the user-established order
    return content.map(item => {
      if (typeof item === 'string') {
        // special divider functionality
        if (item.match(/^-+$/)) return item
        // if we have a string, that's a terminal page. we match it with
        // middleman's page data and return the enhanced object
        return pageData.filter(page => {
          const pageName = page.path
            .split('/')
            .pop()
            .replace(/\.html$/, '')
          return pageName === item
        })[0]
      } else {
        // grab the index page, as it can contain data about the top level link
        item.indexData = pageData.find((page) => {
          const split = page.path.split('/')
          return split[0] === item.category && split[1] === 'index.html'
        })
        // otherwise, it's a nested category. if the category has content, we
        // recurse, passing in that category's content, and the matching
        // subsection of page data from middleman
        if (item.content) {
          item.content = this.matchOrderToPageData(
            item.content,
            this.filterData(pageData, item.category)
          )
        }
        return item
      }
    })
  }

  // recursive render for a recursive data structure!
  renderNavTree(category, content, currentPath) {
    return content.map(item => {
      // dividers are the only items left as strings
      if (typeof item === 'string') return <hr />

      if (item.path) {
        const fileName = item.path.split('/').pop()
        return (
          <li class={currentPath[1] === fileName ? 'active' : ''}>
            <a href={`/${category}/${item.path}`}>
              {item.data.sidebar_title ||
                item.data.page_title ||
                '(!) Page Missing'}
            </a>
          </li>
        )
      } else {
        const current = currentPath[0]
        currentPath.slice(1)
        const title = item.indexData ? (item.indexData.data.sidebar_title || item.indexData.data.page_title) : item.category
        return (
          <li class={current === item.category ? 'active' : ''}>
            <a href={`/${category}/${item.category}/index.html`}>{title}</a>
            {item.content && (
              <ul class="nav">
                {this.renderNavTree(category, item.content, currentPath)}
              </ul>
            )}
          </li>
        )
      }
    })
  }

  filterData(data, category) {
    return data.filter(d => d.path.split('/').includes(category))
  }
}
