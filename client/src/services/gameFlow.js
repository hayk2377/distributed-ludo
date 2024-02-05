import { errorIfUndefinedValues } from '../utils'

export default class GameFlowService {
  constructor({
    gameHost,
    gamePort,
    jwt,

    onConnect,
    onDisconnect,
    onFatalError,
    onNonFatalError,
    onGameStateUpdate,
    onLobbyPlayersUpdate,
  }) {
    errorIfUndefinedValues({
      gameHost,
      gamePort,
      jwt,
    })

    //basic settings
    this.gameHost = gameHost
    this.gamePort = gamePort
    this.jwt = jwt

    //callbacks
    this.updateListeners({
      onConnect,
      onDisconnect,
      onFatalError,
      onNonFatalError,
      onGameStateUpdate,
      onLobbyPlayersUpdate,
    })
  }

  updateListeners = ({
    onConnect,
    onDisconnect,
    onFatalError,
    onNonFatalError,
    onGameStateUpdate,
    onLobbyPlayersUpdate,
  }) => {
    errorIfUndefinedValues({
      onConnect,
      onDisconnect,
      onFatalError,
      onNonFatalError,
      onGameStateUpdate,
      onLobbyPlayersUpdate,
    })
    // server related callbacks
    this.onConnect = onConnect
    this.onDisconnect = onDisconnect
    this.onFatalError = onFatalError
    this.onNonFatalError = onNonFatalError
    this.onGameStateUpdate = onGameStateUpdate
    this.onLobbyPlayersUpdate = onLobbyPlayersUpdate
  }
  createLobby = async () => {
    const res = await fetch(`http://localhost:${this.gamePort}/lobbies`)
    if (!res.ok) throw new Error('Failed to create lobby' + res.statusText)
    const { gameCode } = await res.json()
    return gameCode
  }

  //Either join lobby or rejoin game. intent is either lobby or game
  #join = async ({ gameCode, intent }) => {
    errorIfUndefinedValues({ gameCode, intent })

    const socket = await new Promise((resolve, reject) => {
      const url = `ws://${this.gameHost}:${this.gamePort}/ws?jwt=${this.jwt}&gameCode=${gameCode}&intent=${intent}`
      const socket = new WebSocket(url)

      socket.onclose = () => {
        console.log('Could not connect, check your network connection')
        //If on close with out onmessage, then its a network error
        reject({ message: 'Could not connect, check your network connection', isNetwork: true })
      }
      socket.onopen = () => {
        console.log('Connected to ws, waiting for error or ok')
      }

      socket.onmessage = (e) => {
        const { event, content } = JSON.parse(e.data)
        if (event === 'error') {
          const { isFatal, message } = JSON.parse(content)
          console.error({ isFatal, message })
          reject({ message, isFatal, isNetwork: false })
          return
        }

        if (event === 'ok') {
          console.log('Joined lobby or rejoin game')
          resolve(socket)
          return
        }

        throw new Error(`unknown event ${event}, content is${content}`)
      }
    })

    //set listeners
    this.socket = socket
    this.#setListeners()
    return socket
  }

  joinLobby = async ({ gameCode }) => {
    return await this.#join({ gameCode, intent: 'lobby' })
  }

  rejoinGame = async ({ gameCode }) => {
    return await this.#join({ gameCode, intent: 'game' })
  }

  #getSocket = () => {
    if (!this.socket) throw new Error('Socket doesnt exist')
    return this.socket
  }

  #setListeners = () => {
    const socket = this.#getSocket()
    socket.onopen = () => {
      console.log('connected to ws')
      this.onConnect()
    }

    socket.onclose = () => {
      console.log('disconnected from ws')
      this.onDisconnect()
    }

    socket.onmessage = (e) => {
      const { event, content: contentAsString } = JSON.parse(e.data)
      const content = JSON.parse(contentAsString)

      if (event === 'error') {
        if (content.fatal) {
          console.error(content.message)
          this.onFatalError(content.message)
        } else {
          console.error(content.message)
          this.onNonFatalError(content.message)
        }
      } else if (event === 'gameState') {
        console.log('game state', content)
        this.onGameStateUpdate(content)
      } else if (event === 'lobbyPlayers') {
        console.log('lobby players', content)
        this.onLobbyPlayersUpdate(content)
      } else {
        throw new Error(`unknown event ${event}, content is${content}`)
      }
    }
  }

  #prepareEvent({ event, content }) {
    return JSON.stringify({
      event,
      content: content ? JSON.stringify(content) : '',
    })
  }

  movePawn = (pawnNumber) => {
    console.log('sending move pawn event', pawnNumber)
    const socket = this.#getSocket()
    socket.send(this.#prepareEvent({ event: 'pawn', content: { pawnNumber } }))
  }

  rollDice = () => {
    console.log('sending roll dice event')
    const socket = this.#getSocket()
    socket.send(this.#prepareEvent({ event: 'dice' }))
  }

  startGame = () => {
    console.log('sending start game event')
    const socket = this.#getSocket()
    socket.send(this.#prepareEvent({ event: 'startGame' }))
  }

  disconnect = () => {
    console.log('disconnecting from ws')
    const socket = this.#getSocket()
    if (socket) socket.close()
  }

  endService = ()=>{
    //update listeners to dummies
    this.updateListeners({
      onConnect:()=>{},
      onDisconnect:()=>{},
      onFatalError:()=>{},
      onNonFatalError:()=>{},
      onGameStateUpdate:()=>{},
      onLobbyPlayersUpdate:()=>{},
    })

    //set listeners to these dummy fns
    this.#setListeners()

    //disconnect if there is socket
    if (this.socket) this.disconnect()
  }

  #retryJoining = async ({ gameCode, intent, seconds }) => {
    if (intent !== 'lobby' && intent !== 'game')
      throw new Error('Invalid intent')
    errorIfUndefinedValues({ gameCode, seconds })

    console.log('reconnecting to', intent)

    return new Promise((resolve, reject) => {
      let isPending = true

      const tryConnecting = async () => {
        console.log('trying to reconnect (in tryConnectin()')
        try {
          //Either join lobby or rejoin game (intent=lobby or game)
          this.#join({ gameCode, intent })
        } catch (e) {
          if (e.isNetwork && isPending) {
            //If network error, then try again
            setTimeout(tryConnecting, 1000)
          } else {
            //If prolly logical error, then reject
            reject(e)
            isPending = false
          }
          console.log('reconnect ping failed')
        }
      }

      //Try connecting, if disconnected then tryConnecting will be called again
      tryConnecting()

      //Timeout vs connection. who ever runs first wins the promise
      setTimeout(() => {
        if (isPending) reject(new Error('Reconnect timeout, could not reconnect'))
        isPending = false
      }, seconds * 1000)
    })
  }

  retryJoiningLobby = async ({ gameCode, seconds }) => {
    errorIfUndefinedValues({ gameCode, seconds })
    return await this.#retryJoining({ gameCode, seconds, intent: 'lobby' })
  }

  retryRejoiningGame = async ({ gameCode, seconds }) => {
    errorIfUndefinedValues({ gameCode, seconds })
    return await this.#retryJoining({ gameCode, seconds, intent: 'game' })
  }
}
