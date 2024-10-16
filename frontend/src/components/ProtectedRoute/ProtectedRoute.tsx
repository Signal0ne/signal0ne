import type { ReactNode } from 'react';
import { Navigate } from 'react-router-dom';
import { useAuthContext } from '../../hooks/useAuthContext';

interface ProtectedRouteProps {
  children: ReactNode;
  isDisabled?: boolean;
  redirectTo?: string;
  unAuthed: boolean;
}

const ProtectedRoute = ({
  children,
  isDisabled,
  redirectTo = '/login',
  unAuthed
}: ProtectedRouteProps) => {
  const { currentUser } = useAuthContext();

  if (isDisabled) return <Navigate to={redirectTo} />;

  if (unAuthed) return !currentUser ? children : <Navigate to={redirectTo} />;

  return currentUser ? children : <Navigate to={redirectTo} />;
};

export default ProtectedRoute;
