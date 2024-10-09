import { Link } from 'react-router-dom';
import { SubmitHandler, useForm } from 'react-hook-form';
import { toast } from 'react-toastify';
import { User } from '../../contexts/AuthProvider/AuthProvider';
import { useAuthContext } from '../../hooks/useAuthContext';
import { useState } from 'react';
import Button from '../../components/Button/Button';
import Input from '../../components/Input/Input';
import Spinner from '../../components/Spinner/Spinner';
import './SignUpPage.scss';

interface SignUpFormData {
  confirmPassword: string;
  password: string;
  username: string;
}

interface SignUpResponse {
  accessToken: string;
  refreshToken: string;
  user: User;
}

const SignUpPage = () => {
  const [isSubmitting, setIsSubmitting] = useState(false);

  const {
    formState: { errors },
    handleSubmit,
    register,
    watch
  } = useForm<SignUpFormData>();

  const password = watch('password');

  const { setAccessToken, setCurrentUser } = useAuthContext();

  const handleSignUp: SubmitHandler<SignUpFormData> = async data => {
    const { password, username } = data;

    try {
      setIsSubmitting(true);

      const response = await fetch(
        `${import.meta.env.VITE_SERVER_API_URL}/auth/register`,
        {
          body: JSON.stringify({ password, username }),
          headers: {
            'Content-Type': 'application/json'
          },
          method: 'POST'
        }
      );

      if (!response.ok) throw new Error('Failed to register');

      const data: SignUpResponse = await response.json();

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
    <div className="signup-page">
      <form className="form-container" onSubmit={handleSubmit(handleSignUp)}>
        <div className="form-header">
          <h3 className="form-title">Register 👤</h3>
          <h4 className="form-subtitle">Create your Signal0ne account</h4>
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
              autoComplete="new-password"
              label="Password"
              placeholder="Your password..."
              type="password"
              {...register('password', {
                required: 'This field is required',
                minLength: {
                  message: 'Password must be at least 8 characters long',
                  value: 8
                },
                pattern: {
                  message:
                    'Password must contain at least one letter and one number',
                  value:
                    /^(?=.*[A-Za-z])(?=.*\d)(?=.*[!@#$%^&*])[A-Za-z\d!@#$%^&*]{8,}$/
                }
              })}
            />
            {errors.password && (
              <span className="error-msg">{errors.password?.message}</span>
            )}
          </div>
          <div className="form-field">
            <Input
              autoComplete="new-password"
              label="Confirm Password"
              placeholder="Your password..."
              type="password"
              {...register('confirmPassword', {
                required: 'This field is required',
                validate: value =>
                  value === password || 'Passwords do not match'
              })}
            />
            {errors.password && (
              <span className="error-msg">
                {errors.confirmPassword?.message}
              </span>
            )}
          </div>
        </div>
        <Button
          className="form-submit-btn"
          disabled={isSubmitting}
          type="submit"
        >
          {isSubmitting ? <Spinner /> : 'Register'}
        </Button>
        <p className="form-register-info">
          Already have an account?{' '}
          <strong>
            <Link className="form-register-link" to="/login">
              Login
            </Link>
          </strong>
        </p>
      </form>
    </div>
  );
};

export default SignUpPage;
