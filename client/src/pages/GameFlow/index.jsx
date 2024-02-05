import { useEffect, useState } from 'react'
import CreateJoin from '../CreateJoin'
import Game from '../Game'
import Lobby from '../Lobby'
import GameFlowService from '../../services/gameFlow'
import { ToastContainer, toast } from 'react-toastify'
import 'react-toastify/dist/ReactToastify.css'
import InNavbar from '../../components/InNavbar'

export default function GameFlow() {
  const [page, setPage] = useState(0)
  const [gameCode, setGameCode] = useState(null)
  const [lobbyState, setLobbyState] = useState(null)
  const [hadFatalError, setHadFatalError] = useState(false)
  const [gameState, setGameState] = useState(null)
  const [gameFlowService, setGameFlowService] = useState(null)

  const getPlayerId = () => {
    const urlParams = new URLSearchParams(window.location.search)
    return urlParams.get('id') || '1'
  }

  const getPlayerName = (id) => {
    return { 1: 'Abebe Kebede', 2: 'Kebede Abebe', 3: 'Maritue Mamo' }[id]
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
        toast.error('Disconnected from lobby, Trying to reconnect')
        await gameFlowService.retryJoiningLobby({ gameCode, seconds: 5 })
      } else if (inGame) {
        toast.error('Disconnected from game, Trying to reconnect')
        await gameFlowService.retryRejoiningGame({ gameCode, seconds: 5 })
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

  const createGameFlowService = () => {
    const gameHost = 'localhost'
    const gamePort = 5007
    const jwt = JSON.stringify({
      id: getPlayerId(),
      name: getPlayerName(getPlayerId()),
    })

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

  useEffect(() => {
    setGameFlowService(createGameFlowService())
  }, [])

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
      gameCode = await gameFlowService.createLobby()
    } catch (err) {
      toast.error(`Couldn't connect. Check your connection!`)
      console.error(err)
    }

    if (!gameCode) return
    await gameFlowService.joinLobby({ gameCode })
    //onLobbyPlayersUpdate will handle the page transition
    setGameCode(gameCode)
  }

  const joinLobby = async (gameCode) => {
    try {
      await gameFlowService.joinLobby({ gameCode })
      //onLobbyPlayersUpdate will handle the page transition
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

  const toHomeOnGameEnd = () => {
    //Just like the use effect
    gameFlowService.endService()
    setGameFlowService(createGameFlowService())
    setPage(0)
  }

  return (
    <div className='flex flex-col gap-6'>
      <InNavbar />
      <main>
        <>
          <ToastContainer />
          {page === 0 && (
            <CreateJoin onCreateLobby={createLobby} onJoinLobby={joinLobby} />
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
