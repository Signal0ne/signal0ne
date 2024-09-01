import './LoginPage.scss';

const LoginPage = () => {
  return (
    <div className="login-page">
      <form className="form-container" onSubmit={e => e.preventDefault()}>
        <div className="form-header">
          <h3 className="form-title">Log In</h3>
          <h4 className="form-subtitle">with Signal0ne account</h4>
        </div>
        <div className="form-content">
          <div className="form-field">
            <label className="form-label" htmlFor="email">
              Email
            </label>
            <input
              className="form-input"
              id="email"
              placeholder="Email"
              type="email"
            />
          </div>
          <div className="form-field">
            <label className="form-label" htmlFor="password">
              Password
            </label>
            <input
              className="form-input"
              id="password"
              placeholder="Password"
              type="password"
            />
          </div>
        </div>
        <button className="form-submit-btn" type="submit">
          Log In
        </button>

        <hr className="form-separator" />
      </form>
    </div>
  );
};

export default LoginPage;
