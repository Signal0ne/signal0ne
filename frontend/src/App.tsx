import { Tooltip } from 'react-tooltip';
import AppRoutes from './components/AppRoutes/AppRoutes';
import Navbar from './components/Navbar/Navbar';
import './App.scss';

const App = () => {
  return (
    <>
      <Navbar />
      <AppRoutes />
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
