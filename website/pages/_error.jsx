import NotFound from './404'
import Bugsnag from '@hashicorp/nextjs-scripts/lib/bugsnag'

function Error({ statusCode }) {
  return <NotFound statusCode={statusCode} />
}

Error.getInitialProps = ({ res, err }) => {
  if (err) Bugsnag.notify(err)
  const statusCode = res ? res.statusCode : err ? err.statusCode : 404
  return { statusCode }
}

export default Error
