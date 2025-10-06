import { Button } from "./components/ui/button"
import { ScrollArea } from "./components/ui/scroll-area"

function App() {
  return (
    <main className='h-screen w-full overflow-hidden text-slate-100 relative'>
      <div className='h-10 w-full sticky top-0 z-10 bg-gray-700/60'>
        <div className='drag-zone cursor-move flex items-center justify-between px-0.5 z-0 font-serif'>
          <h1 className='text-2xl font-bold mb-2 pt-0.5 pl-1'>goDrawer</h1>
          <Button className='bg-orange-500 hover:bg-orange-300'>
            Settings
          </Button>
        </div>
      </div>
      <ScrollArea className='h-44 px-2 text-base'>
        <p>Drag files here to open them asdfas dfas </p>
        <p>Drag files he一些中午</p>
        <p>Drag files here to open them</p>
        <p>Drag files here to open them</p>
        <p>Drag files here to open them</p>
        <p>Drag files here to open them</p>
        <p>Drag files here to open them</p>
        <p>Drag files here to open them</p>
        <p>Drag files here to open them</p>
        <p>Drag files here to open them</p>
        <p>Drag files here to open them</p>
        <p>Drag files here to open them</p>
        <p>Drag files here to open them</p>
      </ScrollArea>
      <div className='flex items-center justify-center pt-0.5'>
        <Button className='bg-sky-500 hover:bg-sky-300 font-serif'>
          Add Drawer
        </Button>
      </div>
      {/* </div> */}
    </main>
  )
}

export default App
