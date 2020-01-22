import React from 'react'
import bugsnag from '@bugsnag/js'
import bugsnagReact from '@bugsnag/plugin-react'

const apiKey =
  typeof window === 'undefined'
    ? 'fb2dc40bb48b17140628754eac6c1b11'
    : '07ff2d76ce27aded8833bf4804b73350'

const bugsnagClient = bugsnag({
  apiKey,
  releaseStage: process.env.NODE_ENV || 'development'
})

bugsnagClient.use(bugsnagReact, React)

export default bugsnagClient
