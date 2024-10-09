import { Link } from 'react-router-dom';
import { SubmitHandler, useForm } from 'react-hook-form';
import { toast } from 'react-toastify';
import { User } from '../../contexts/AuthProvider/AuthProvider';
import { useAuthContext } from '../../hooks/useAuthContext';
import { useState } from 'react';
import Button from '../../components/Button/Button';
import Input from '../../components/Input/Input';
import Spinner from '../../components/Spinner/Spinner';
import './LoginPage.scss';

interface LoginFormData {
  password: string;
  username: string;
}

interface LoginResponse {
  accessToken: string;
  refreshToken: string;
  user: User;
}

const LoginPage = () => {
  const [isSubmitting, setIsSubmitting] = useState(false);

  const { setAccessToken, setCurrentUser } = useAuthContext();

  const {
    formState: { errors },
    handleSubmit,
    register
  } = useForm<LoginFormData>();

  const handleLogin: SubmitHandler<LoginFormData> = async data => {
    const { password, username } = data;

    try {
      setIsSubmitting(true);

      const response = await fetch(
        `${import.meta.env.VITE_SERVER_API_URL}/auth/login`,
        {
          body: JSON.stringify({ password, username }),
          headers: {
            'Content-Type': 'application/json'
          },
          method: 'POST',
          credentials: 'include'
        }
      );

      if (!response.ok) throw new Error('Failed to login');

      const data: LoginResponse = await response.json();

      setAccessToken(data.accessToken);
      setCurrentUser(data.user);
      saveToLocalStorage(data.user);
    } catch (error) {
      if (error instanceof Error) {
        toast.error(error.message);
      } else {
        toast.error('An unknown error occurred. Please try again later.');
      }
    } finally {
      setIsSubmitting(false);
    }
  };

  const saveToLocalStorage = (data: User) => {
    localStorage.setItem('user', JSON.stringify(data));
  };

  return (
    <div className="login-page">
      <form className="form-container" onSubmit={handleSubmit(handleLogin)}>
        <div className="form-header">
          <h3 className="form-title">Welcome Back 👋</h3>
          <h4 className="form-subtitle">Login to your Signal0ne account</h4>
        </div>
        <div className="form-content">
          <div className="form-field">
            <Input
              autoComplete="username"
              label="Username"
              placeholder="Your username..."
              {...register('username', {
                required: 'This field is required'
              })}
            />
            {errors.username && (
              <span className="error-msg">{errors.username?.message}</span>
            )}
          </div>
          <div className="form-field">
            <Input
              autoComplete="current-password"
              label="Password"
              placeholder="Your password..."
              type="password"
              {...register('password', {
                required: 'This field is required'
              })}
            />
            {errors.password && (
              <span className="error-msg">{errors.password?.message}</span>
            )}
          </div>
        </div>
        <Button
          className="form-submit-btn"
          disabled={isSubmitting}
          type="submit"
        >
          {isSubmitting ? <Spinner /> : 'Log In'}
        </Button>
        <p className="form-register-info">
          Don't have an account?{' '}
          <strong>
            <Link className="form-register-link" to="/register">
              Register
            </Link>
          </strong>
        </p>
      </form>
    </div>
  );
};

export default LoginPage;
