import { AccountIcon, GearIcon, Signal0neLogo } from '../Icons/Icons';
import { handleKeyDown } from '../../utils/utils';
import { Link } from 'react-router-dom';
import { ROUTES } from '../../data/routes';
import { useAuthContext } from '../../hooks/useAuthContext';
import { useState } from 'react';
import Button from '../Button/Button';
import './Navbar.scss';

const Navbar = () => {
  const [isAccountOpen, setIsAccountOpen] = useState(false);

  const { currentUser, logout } = useAuthContext();

  const handleLogout = async () => {
    await logout();
    setIsAccountOpen(false);
  };

  const handleOpenAccount = () => setIsAccountOpen(prev => !prev);

  const getNavbarLinks = () =>
    currentUser ? (
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
              onKeyDown={handleKeyDown(handleOpenAccount)}
              tabIndex={0}
              width={36}
            />
            {isAccountOpen && (
              <div className="account-content">
                <span className="account-name">
                  User: <strong>{currentUser.name}</strong>
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
