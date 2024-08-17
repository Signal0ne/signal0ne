import { AccountIcon, GearIcon, Signal0neLogo } from '../Icons/Icons';
import { Link } from 'react-router-dom';
import { ROUTES } from '../../data/routes';
import { useAuthContext } from '../../hooks/useAuthContext';
import './Navbar.scss';

const Navbar = () => {
  const { currentUser } = useAuthContext();
  console.log(currentUser);
  return (
    <nav className="navbar-container">
      <div className="navbar-logo-container">
        <Link to="/">
          <Signal0neLogo width="120px" />
        </Link>
      </div>
      <div className="navbar-content-container">
        {currentUser ? (
          <>
            <div className="navbar-content-links">
              {ROUTES.map(({ isDisabled, path, title }) =>
                isDisabled ? (
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
                )
              )}
            </div>
            <div className="navbar-content-actions">
              <GearIcon height={28} tabIndex={0} width={28} />
              <AccountIcon height={32} tabIndex={0} width={32} />
            </div>
          </>
        ) : (
          <div className="navbar-content-links">
            <Link className="navbar-content-link" to="/login" key="login">
              Sign In
            </Link>
          </div>
        )}
      </div>
    </nav>
  );
};

export default Navbar;
