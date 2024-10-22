import { AuthContextProvider } from './contexts/AuthProvider/AuthProvider';
import { QueryClient, QueryClientProvider } from '@tanstack/react-query';
import { ToastContainer } from 'react-toastify';
import { Tooltip } from 'react-tooltip';
import AppRoutes from './components/AppRoutes/AppRoutes';
import Navbar from './components/Navbar/Navbar';
import 'react-toastify/dist/ReactToastify.css';
import './App.scss';

const queryClient = new QueryClient();

const App = () => (
  <QueryClientProvider client={queryClient}>
    <AuthContextProvider>
      <Navbar />
      <AppRoutes />
      <ToastContainer />
    </AuthContextProvider>
    <Tooltip
      delayShow={50}
      id="global"
      openEvents={{ focus: true, mouseover: true }}
      variant="light"
    />
  </QueryClientProvider>
);

export default App;
