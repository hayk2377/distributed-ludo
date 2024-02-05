import { Link, useNavigate } from 'react-router-dom'
import { COLORS } from '../utils'
import { hasJWT, logout } from '../services/user'
import { toast } from 'react-toastify'

export default function InNavbar() {
  const navigate = useNavigate()
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
        <div className='flex gap-3 items-center'>
          <Link to='/app/'>Game</Link>
          <button
            onClick={() => {
              try {
                logout()
                navigate('/')
                console.log('logged out, jwt is', hasJWT())
              } catch (e) {
                toast.error(e.message)
                console.log(e.message)
              }
            }}
          >
            Logout
          </button>
        </div>
      </nav>
    </div>
  )
}
