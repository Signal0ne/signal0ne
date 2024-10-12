import { AuthContextProvider } from './contexts/AuthProvider/AuthProvider';
import { ToastContainer } from 'react-toastify';
import { Tooltip } from 'react-tooltip';
import AppRoutes from './components/AppRoutes/AppRoutes';
import Navbar from './components/Navbar/Navbar';
import 'react-toastify/dist/ReactToastify.css';
import './App.scss';

const App = () => {
  return (
    <>
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
    </>
  );
};

export default App;
