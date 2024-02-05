import OutNavbar from '../../components/OutNavbar'
import { COLORS } from '../../utils'
import { Link, useNavigate } from 'react-router-dom'
import { ToastContainer, toast } from 'react-toastify'
import 'react-toastify/dist/ReactToastify.css'
import { login } from '../../services/user'

export default function Login() {
  const navigate = useNavigate()

  const onLogin = async (e) => {
    e.preventDefault()
    const email = e.target.email.value
    const password = e.target.password.value
    await login({ email, password })
    toast.success('Login successful'+JSON.stringify({email, password}))
    e.preventDefault()
    await new Promise((resolve) => setTimeout(resolve, 1000))
    navigate('/app')
  }
  return (
    <>
      <ToastContainer />
      <div className='flex flex-col gap-6'>
        <OutNavbar />
        <main>
          <div className='container'>
            <div className='flex-1 flex flex-col gap-3'>
              <h1 className='font-bold text-3xl'>Login</h1>
              <form className='flex flex-col gap-3' onSubmit={onLogin}>
                <label>
                  Email
                  <br />
                  <input
                    type='email'
                    name='email'
                    className='border p-3 rounded'
                    placeholder='your@email.com'
                    defaultValue={'abebe@gmail.com'}
                    required
                  />
                </label>
                <label>
                  Password
                  <br />
                  <input
                    type='password'
                    name='password'
                    className='border p-3 rounded'
                    defaultValue={'abebe'}
                    required
                  />
                </label>

                <div>
                  <button
                    className='p-3 rounded text-white'
                    style={{ background: COLORS.green }}
                  >
                    Login
                  </button>
                </div>
              </form>
            </div>
          </div>
        </main>
      </div>
    </>
  )
}
