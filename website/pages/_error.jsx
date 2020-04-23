import React from 'react'
import Bugsnag from '../lib/bugsnag'
import FourOhFour from './404'

export default class Page extends React.Component {
  static async getInitialProps(ctx) {
    if (ctx.err) Bugsnag.notify(ctx.err)
  }

  render() {
    return <FourOhFour />
  }
}
