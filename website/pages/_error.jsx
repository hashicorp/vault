import React from 'react'
import bugsnagClient from '../lib/bugsnag'
import FourOhFour from './404'

export default class Page extends React.Component {
  static async getInitialProps(ctx) {
    if (ctx.err) bugsnagClient.notify(ctx.err)
  }

  render() {
    return <FourOhFour />
  }
}
