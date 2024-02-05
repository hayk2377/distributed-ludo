import OutNavbar from '../../components/OutNavbar'
import { Link } from 'react-router-dom'
import { COLORS } from '../../utils'

export default function LandingPage() {
  return (
    <div className='flex flex-col gap-6'>
      <OutNavbar />
      <main>
        <div className='container'>
          <div className='flex gap-6'>
            <div className='flex-1 flex flex-col gap-3'>
              <h1 className='font-bold text-3xl'>Ludoking plus S?</h1>
              <p>
                Your beloved game Ludoking has now become online! You can now
                continue winning your friendS and family memberS where ever you
                are!
              </p>

              <div className='flex gap-3'>
                <Link to='/signup'>
                  <button
                    className='p-3 rounded text-white'
                    style={{ background: COLORS.green }}
                  >
                    Signup
                  </button>
                </Link>

                <Link to='/login'>
                  <button
                    className='border p-3 rounded'
                    style={{ borderColor: COLORS.green, color: COLORS.green }}
                  >
                    Login
                  </button>
                </Link>
              </div>
            </div>

            <div className='flex-1'></div>
          </div>
        </div>
      </main>
    </div>
  )
}
