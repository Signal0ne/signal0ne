import { Route, Routes } from 'react-router-dom';
import { ROUTES } from '../../data/routes';
import ProtectedRoute from '../ProtectedRoute/ProtectedRoute';

const AppRoutes = () => (
  <Routes>
    {ROUTES.map(({ Component, path, unAuthed }) => (
      <Route
        element={
          <ProtectedRoute unAuthed={unAuthed}>
            <Component />
          </ProtectedRoute>
        }
        key={path}
        path={path}
      />
    ))}
  </Routes>
);

export default AppRoutes;
