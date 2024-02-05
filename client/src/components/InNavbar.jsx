import { Link } from 'react-router-dom'
import { COLORS } from '../utils'

export default function InNavbar() {
  return (
    <div className=' shadow py-3'>
      <nav className='container flex justify-between'>
        <Link to='/' className='font-bold text-xl'>
          Ludo
          <span style={{ color: COLORS.red }}>K</span>
          <span style={{ color: COLORS.green }}>I</span>
          <span style={{ color: COLORS.blue }}>N</span>
          <span style={{ color: COLORS.yellow }}>G</span>S
        </Link>
        <div className='flex gap-3'>
          <Link to='/app/'>Game</Link>
          <Link to='/app/settings'>Settings</Link>
        </div>
      </nav>
    </div>
  )
}
