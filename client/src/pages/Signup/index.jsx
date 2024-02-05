import { useState } from 'react'
import OutNavbar from '../../components/OutNavbar'
import { COLORS } from '../../utils'
import { Link, useNavigate } from 'react-router-dom'
import { ToastContainer, toast } from 'react-toastify'
import 'react-toastify/dist/ReactToastify.css'
import { signup, test } from '../../services/user'

export default function Signup() {
  const navigate = useNavigate()

  const onSignup = async (e) => {
    e.preventDefault()
    const email = e.target.email.value
    const name = e.target.name.value
    const password = e.target.password.value
    try {
      console.log('signingup')
      await signup({ email, name, password })
      // await test()
      console.log('here')
      toast.success('Saved' + JSON.stringify({ name, password, email }))
      await new Promise((resolve) => setTimeout(resolve, 1000))
      navigate('/login')
    } catch (e) {
      return toast.error(e.message)
    }
  }
  return (
    <>
      <ToastContainer />
      <div className='flex flex-col gap-6'>
        <OutNavbar />
        <main>
          <div className='container'>
            <div className='flex-1 flex flex-col gap-3'>
              <h1 className='font-bold text-3xl'>Signup</h1>
              <form className='flex flex-col gap-3' onSubmit={onSignup}>
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
                  Name
                  <br />
                  <input
                    type='text'
                    name='name'
                    className='border p-3 rounded'
                    placeholder='Your Player Name'
                    defaultValue={'Abebe Kebede'}
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
                    Signup
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
