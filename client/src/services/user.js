import { errorIfUndefinedValues } from '../utils'
import { LOAD_BALANCER_ADDRESS as AUTH_ADDRESS } from './loadbalancer-address' 


export function getJWT() {
  const token = localStorage.getItem('token')
  if (!token) throw new Error('No token stored')
  return token
}

export function removeJWT() {
  localStorage.removeItem('token')
}

export function hasJWT() {
  if (localStorage.getItem('token')) return true
  else return false
}

async function extractErrorMessage(res) {
  try {
    const { error } = await res.json()
    return error
  } catch (e) {
    return res.statusText
  }
}

export function setJWT(token) {
  localStorage.setItem('token', token)
}

export async function signup({
  email,
  password,
  name,
  authAddress = AUTH_ADDRESS,
}) {
  errorIfUndefinedValues({ email, password, name })

  console.log('here1')
  const url = `${authAddress}/signup`
  const res = await fetch(url, {
    method: 'POST',
    headers: { 'Content-Type': 'application/json' },
    body: JSON.stringify({ email, password, name }),
  })

  console.log('here2')
  if (!res.ok)
    throw new Error('Failed to signup' + (await extractErrorMessage(res)))
}

export async function login({ email, password, authAddress = AUTH_ADDRESS }) {
  errorIfUndefinedValues({ email, password })
  const url = `${authAddress}/login`
  const res = await fetch(url, {
    method: 'POST',
    headers: { 'Content-Type': 'application/json' },
    body: JSON.stringify({ email, password }),
  })

  if (!res.ok)
    throw new Error('Failed to login' + (await extractErrorMessage(res)))
  const data = await res.json()
  console.log('jwt was', data)
  const { token } = data
  setJWT(token)
}

export async function logout() {
  removeJWT()
}

export async function getUser(authAddress = AUTH_ADDRESS ) {
  if (!hasJWT()) throw new Error('No token stored')

  const url = `${authAddress}/user`
  const res = await fetch(url, {
    headers: {
      'Content-Type': 'application/json',
      Authorization: `Bearer ${getJWT()}`,
    },
  })

  if (!res.ok)
    throw new Error('Failed to get user' + (await extractErrorMessage(res)))

  return await res.json()
}

export async function test() {
  const url = `${AUTH_ADDRESS}/test`
  const res = await fetch(url)

  if (!res.ok)
    throw new Error('Failed to test' + (await extractErrorMessage(res)))

  return await res.json()
}
