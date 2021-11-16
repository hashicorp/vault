// Imports below are used in getStaticProps only
import RAW_CONTENT from './content.json'

export async function getStaticProps() {
  return { props: {} }
}

export default function Homepage({ content }) {
  return <p>Homepage</p>
}
