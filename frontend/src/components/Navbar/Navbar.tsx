import { AccountIcon, GearIcon, Signal0neLogo } from '../Icons/Icons';
import { Link } from 'react-router-dom';
import { ROUTES } from '../../data/routes';
import { User } from '../../contexts/AuthProvider/AuthProvider';
import { useAuthContext } from '../../hooks/useAuthContext';
import { useRef, useState } from 'react';
import Button from '../Button/Button';
import './Navbar.scss';

const Navbar = () => {
  const [isAccountOpen, setIsAccountOpen] = useState(false);

  const { setCurrentUser } = useAuthContext();

  const userRef = useRef<User | null>(null);

  if (userRef.current === null) {
    const userString = localStorage.getItem('user');

    if (userString) {
      userRef.current = JSON.parse(userString);
    }
  }

  const handleLogout = async () => {
    try {
      const response = await fetch(
        `${import.meta.env.VITE_SERVER_API_URL}/auth/logout`,
        {
          credentials: 'include'
        }
      );

      if (!response.ok) throw new Error('Failed to logout');

      setCurrentUser(null);
      localStorage.removeItem('user');
      userRef.current = null;
    } catch (error) {
      console.error(error);
    }
  };

  const handleOpenAccount = () => setIsAccountOpen(prev => !prev);

  const getNavbarLinks = () =>
    userRef.current ? (
      <>
        <div className="navbar-content-links">
          {ROUTES.map(({ isDisabled, path, showInNavbar, title, unAuthed }) => {
            if (unAuthed || !showInNavbar) return;

            return isDisabled ? (
              <span
                className="navbar-content-link disabled"
                data-tooltip-content="Coming soon..."
                data-tooltip-id="global"
                key={path}
                tabIndex={0}
              >
                {title}
              </span>
            ) : (
              <Link className="navbar-content-link" to={path} key={path}>
                {title}
              </Link>
            );
          })}
        </div>
        <div className="navbar-content-actions">
          <GearIcon height={32} tabIndex={0} width={32} />
          <div className="account-container">
            <AccountIcon
              height={36}
              onClick={handleOpenAccount}
              tabIndex={0}
              width={36}
            />
            {isAccountOpen && (
              <div className="account-content">
                <span className="account-name">
                  User: <strong>{userRef.current.name}</strong>
                </span>
                <hr className="separator" />
                <Button className="account-logout-btn" onClick={handleLogout}>
                  Log Out
                </Button>
              </div>
            )}
          </div>
        </div>
      </>
    ) : (
      <>
        <div className="navbar-content-links">
          <Link className="navbar-content-link" to="/login">
            Sign In
          </Link>
          <Link className="navbar-content-link cta" to="/register">
            Sign Up
          </Link>
        </div>
      </>
    );

  return (
    <nav className="navbar-container">
      <div className="navbar-logo-container">
        <Link to="/">
          <Signal0neLogo width="120px" />
        </Link>
      </div>
      <div className="navbar-content-container">{getNavbarLinks()}</div>
    </nav>
  );
};

export default Navbar;
