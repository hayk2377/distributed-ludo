import { useEffect, useState } from 'react'
import InNavbar from '../../components/InNavbar'
import { COLORS } from '../../utils'
import { Link, useNavigate } from 'react-router-dom'
import { ToastContainer, toast } from 'react-toastify'
import 'react-toastify/dist/ReactToastify.css'

export default function Settings() {
  const navigate = useNavigate()

  const [email, setEmail] = useState('')
  const [name, setName] = useState('')
  const [password, setPassword] = useState('')

  useEffect(() => {
    const fetchUser = async () => {
      const user = {email:'abebe@gmail.com', name:'Abebe Kebede', password:'abebe'}
      setEmail(user.email)
      setName(user.name)
      setPassword(user.password)
    }
    fetchUser()
  }, [])

  const onSave = async (e) => {
    e.preventDefault()
    const name = e.target.name.value
    const password = e.target.password.value
    toast.success('Saved' + JSON.stringify({name, password}))
    // await new Promise((resolve) => setTimeout(resolve, 1000))
  }
  return (
    <>
      <ToastContainer />
      <div className='flex flex-col gap-6'>
        <InNavbar />
        <main>
          <div className='container'>
            <div className='flex-1 flex flex-col gap-3'>
              <h1 className='font-bold text-3xl'>Settings</h1>
              <form className='flex flex-col gap-3' onSubmit={onSave}>
                <label className='opacity-50'>
                  Email
                  <br />
                  <input
                    type='email'
                    name='email'
                    disabled
                    className=' border p-3 rounded'
                    placeholder='your@email.com'
                    defaultValue={email}
                  />
                </label>
                <label>
                  Name
                  <br />
                  <input
                    type='text'
                    name='name'
                    className='border p-3 rounded'
                    placeholder='Your Player Name'
                    defaultValue={name}
                  />
                </label>
                <label>
                  Password
                  <br />
                  <input
                    type='password'
                    name='password'
                    className='border p-3 rounded'
                    defaultValue={password}
                  />
                </label>

                <div>
                  <button
                    className='p-3 rounded text-white'
                    style={{ background: COLORS.green }}
                  >
                    Save
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

