import { Navigate } from 'react-router-dom';
import { ReactNode } from 'react';
import { useAuthContext } from '../../hooks/useAuthContext';

interface ProtectedRouteProps {
  children: ReactNode;
  redirectTo?: string;
  unAuthed: boolean;
}

const ProtectedRoute = ({
  children,
  redirectTo = '/login',
  unAuthed
}: ProtectedRouteProps) => {
  const { currentUser } = useAuthContext();

  return children;

  //TODO: Uncomment after finishing the authentication logic
  if (unAuthed) return !currentUser ? children : <Navigate to={redirectTo} />;

  return currentUser ? children : <Navigate to={redirectTo} />;
};

export default ProtectedRoute;
