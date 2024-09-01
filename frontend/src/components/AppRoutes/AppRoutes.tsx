import { Route, Routes } from 'react-router-dom';
import { ROUTES } from '../../data/routes';
import ProtectedRoute from '../ProtectedRoute/ProtectedRoute';
import NotFoundPage from '../../pages/NotFoundPage/NotFoundPage';

const AppRoutes = () => (
  <Routes>
    {ROUTES.map(({ Component, path, redirectTo, unAuthed }) => (
      <Route
        element={
          <ProtectedRoute redirectTo={redirectTo} unAuthed={unAuthed}>
            <Component />
          </ProtectedRoute>
        }
        key={path}
        path={path}
      />
    ))}
    <Route element={<NotFoundPage />} path="*" />
  </Routes>
);

export default AppRoutes;
