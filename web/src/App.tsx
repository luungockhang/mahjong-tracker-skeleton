
// web/src/App.tsx
import { useQuery } from '@tanstack/react-query'

function App() {
  const { data, isLoading } = useQuery({
    queryKey: ['health'],
    queryFn: async () => {
      const res = await fetch('/api/healthz')
      return res.json()
    },
  })

  if (isLoading) return <p>Loading...</p>

  return (
    <div>
      <h1>Mahjong Tracker</h1>
      <pre>{JSON.stringify(data)}</pre>
    </div>
  )
}

export default App
