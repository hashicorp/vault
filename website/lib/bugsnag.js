import React from 'react'
import Bugsnag from '@bugsnag/js'
import BugsnagReact from '@bugsnag/plugin-react'

const apiKey =
  typeof window === 'undefined'
    ? 'fb2dc40bb48b17140628754eac6c1b11' // server key
    : '07ff2d76ce27aded8833bf4804b73350' // client key

if (!Bugsnag._client) {
  Bugsnag.start({
    apiKey,
    plugins: [new BugsnagReact(React)],
    otherOptions: { releaseStage: process.env.NODE_ENV || 'development' },
  })
}

export default Bugsnag
