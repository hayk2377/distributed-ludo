import { useEffect, useState } from 'react'
import CreateJoin from '../CreateJoin'
import Game from '../Game'
import Lobby from '../Lobby'
import GameFlowService, {
  createLobby as cLobby,
  getGameServerAddress,
} from '../../services/gameFlow'
import { ToastContainer, toast } from 'react-toastify'
import 'react-toastify/dist/ReactToastify.css'
import InNavbar from '../../components/InNavbar'
import { getJWT, getUser } from '../../services/user'

export default function GameFlow() {
  const [user, setUser] = useState(null)
  const [page, setPage] = useState(0)
  const [gameCode, setGameCode] = useState(null)
  const [lobbyState, setLobbyState] = useState(null)
  const [hadFatalError, setHadFatalError] = useState(false)
  const [gameState, setGameState] = useState(null)
  const [gameFlowService, setGameFlowService] = useState(null)

  //First of all, say welcom
  useEffect(() => {
    const fetchUser = async () => {
      try {
        const user = await getUser()
        setUser(user)
      } catch (e) {
        toast.error(e.message)
        toast.error('Please retry')
      }
    }
    fetchUser()
  }, [])

  const getPlayerId = () => {
    return user.id
  }

  const getPlayerName = (id) => {
    return user.name
  }

  const onConnect = () => console.log('Connected')
  const onDisconnect = async () => {
    const elseWhere = page === 0
    const inLobby = page === 1
    const inGame = page === 2

    console.log(
      `disconnect in elseWher: ${elseWhere} in lobby: ${inLobby} in game: ${inGame}`
    )

    if (hadFatalError) return

    if (elseWhere) {
      toast.error("Couldn't connect. Check your connection!")
      console.log('Join prolly is not working')
      return
    }

    try {
      if (inLobby) {
        toast.error('Disconnected from lobby, Trying to reconnect to server')
        await gameFlowService.retryJoiningLobby({ gameCode, seconds: 5 })
      } else if (inGame) {
        toast.error('Disconnected from game, Trying to reconnect')
        // await gameFlowService.retryRejoiningGame({ gameCode, seconds: 5 })
        for (let i = 0; i < 5; i++) {
          await new Promise((resolve) => setTimeout(resolve, 1000))
          try {
            toast.warn('trying to reconnect via load balancer')
            const gameFlowService = await createGameFlowService(gameCode)
            await gameFlowService.rejoinGame({ gameCode })
            setGameFlowService(gameFlowService)
            toast.success('reconnection worked!')
            break
          } catch (e) {
            toast.error('reconnection error: ', e)
          }
        }
      }
    } catch (err) {
      toast.error(err.message)
      console.error('reconnection error: ', err)
      setPage(0)
    }
    console.log('Disconnected')
  }

  const onFatalError = (error) => {
    toast.error(error)
    setHadFatalError(true)
    console.error('Fatal error:', error)
  }
  const onNonFatalError = (error) => {
    toast.error(error)
    console.error('Error:', error)
  }
  const onGameStateUpdate = (gameState) => {
    setGameState(gameState)
    console.log('Game state:', gameState)
    setPage(2)
  }

  const onLobbyPlayersUpdate = (lobbyState) => {
    setPage(1)
    setLobbyState(lobbyState)
    console.log('Lobby state:', lobbyState)
  }

  const createGameFlowService = async (gameCode) => {
    if (!gameCode) throw new Error('Game code is required')
    const { host: gameHost, port: gamePort } = await getGameServerAddress(
      gameCode
    )

    const jwt = getJWT()

    const newGameFlowService = new GameFlowService({
      gameHost,
      gamePort,
      jwt,
      onConnect,
      onDisconnect,
      onFatalError,
      onNonFatalError,
      onGameStateUpdate,
      onLobbyPlayersUpdate,
    })

    return newGameFlowService
  }

  // useEffect(() => {
  //   const fetch = async () => {
  //     try {
  //       const gameFlowService = await createGameFlowService()
  //       setGameFlowService(gameFlowService)
  //     } catch (e) {
  //       toast.error(e.message)
  //       toast.error('Please retry')
  //     }
  //   }
  //   fetch()
  // }, [])

  useEffect(() => {
    //Some listeners depend on page
    gameFlowService?.updateListeners({
      onConnect,
      onDisconnect,
      onFatalError,
      onNonFatalError,
      onGameStateUpdate,
      onLobbyPlayersUpdate,
    })
  }, [page])

  const createLobby = async () => {
    console.log('Creating lobby')
    let gameCode
    try {
      gameCode = await cLobby()
      console.log('game code is', gameCode)
      const gameFlowService = await createGameFlowService(gameCode)
      await gameFlowService.joinLobby({ gameCode })

      setGameCode(gameCode)
      setGameFlowService(gameFlowService)
    } catch (err) {
      toast.error(`Couldn't connect. Check your connection!`)
      console.error(err)
    }
  }

  const joinLobby = async (gameCode) => {
    try {
      const gameFlowService = await createGameFlowService(gameCode)
      await gameFlowService.joinLobby({ gameCode })
      //onLobbyPlayersUpdate will handle the page transition

      setGameCode(gameCode)
      setGameFlowService(gameFlowService)
    } catch (err) {
      toast.error(err.message)
      console.error(err)
    }
  }

  const startGame = async () => {
    try {
      await gameFlowService.startGame()
    } catch (err) {
      //onGameStateUpdate will handle the page transition
      toast.error(err.message)
      console.error(err)
    }
  }

  const movePawn = (pawnNumber) => {
    gameFlowService.movePawn(pawnNumber)
  }
  const rollDice = () => {
    gameFlowService.rollDice()
  }

  const toHomeOnGameEnd = async () => {
    //Just like the use effect
    gameFlowService.endService()
    const gF = await createGameFlowService()
    setGameFlowService(gF)
    setPage(0)
  }

  return (
    <div className='flex flex-col gap-6'>
      <InNavbar />
      <main>
        <>
          <ToastContainer />
          {page === 0 && (
            <CreateJoin
              onCreateLobby={createLobby}
              onJoinLobby={joinLobby}
              name={user?.name}
            />
          )}
          {page === 1 && (
            <Lobby
              lobbyState={lobbyState}
              onStartGame={startGame}
              playerId={getPlayerId()}
            />
          )}
          {page === 2 && (
            <Game
              gameState={gameState}
              onPawn={movePawn}
              onRollDice={rollDice}
              playerId={getPlayerId()}
              onGameEnd={toHomeOnGameEnd}
            />
          )}
        </>
      </main>
    </div>
  )
}
