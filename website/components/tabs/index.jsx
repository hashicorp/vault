import ReactTabs from '@hashicorp/react-tabs'

export default function Tabs({ children }) {
  return (
    <ReactTabs
      items={children.map((Block) => ({
        heading: Block.props.heading,
        // eslint-disable-next-line react/display-name
        tabChildren: () => Block,
      }))}
    />
  )
}

export function Tab({ children }) {
  return <>{children}</>
}
