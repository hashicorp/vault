import s from './style.module.css'

export default function Columns({ count = 1, children }) {
  return <div className={`${s.root} ${s['count-' + count]}`}>{children}</div>
}
