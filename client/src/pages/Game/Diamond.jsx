import { COLORS } from '../../utils'

export default function Diamond() {
  const colors = []

  return (
    <div className='absolute w-[20%] h-[20%] top-[50%] left-[50%] translate-x-[-50%] translate-y-[-50%]'>
      <div
        className='relative tile-border w-full h-full overflow-hidden'
      >
        <div
          className='w-full h-full absolute top-[0%] left-[50%] rotate-45 translate-x-[-50%] translate-y-[-50%] scale-[0.707]'
          style={{ background: COLORS.green }}
        ></div>
        <div
          className='w-full h-full absolute top-[50%] left-[0%] rotate-45 translate-x-[-50%] translate-y-[-50%] scale-[0.707]'
          style={{ background: COLORS.red }}
        ></div>
        <div
          className='w-full h-full absolute top-[100%] left-[50%]  rotate-45 translate-x-[-50%] translate-y-[-50%] scale-[0.707]'
          style={{ background: COLORS.blue }}
        ></div>
        <div
          className='w-full h-full absolute top-[50%] left-[100%] rotate-45 translate-x-[-50%] translate-y-[-50%] scale-[0.707]'
          style={{ background: COLORS.yellow }}
        ></div>
      </div>
    </div>
  )
}
