import React from 'react'
import ReactDOM from 'react-dom/client'
import { QueryClient, QueryClientProvider } from '@tanstack/react-query'
import { MantineProvider } from '@mantine/core'
import App from './App'


const queryClient = new QueryClient()


ReactDOM.createRoot(document.getElementById('root')!).render(
<React.StrictMode>
<MantineProvider>
<QueryClientProvider client={queryClient}>
<App />
</QueryClientProvider>
</MantineProvider>
</React.StrictMode>
)