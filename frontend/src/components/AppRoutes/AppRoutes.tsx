import { Route, Routes } from 'react-router-dom';
import { ROUTES } from '../../data/routes';

const AppRoutes = () => (
  <Routes>
    {ROUTES.map(({ Component, path }) => (
      <Route element={<Component />} key={path} path={path} />
    ))}
  </Routes>
);

export default AppRoutes;
