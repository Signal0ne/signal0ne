import { Link } from 'react-router-dom';
import { SubmitHandler, useForm } from 'react-hook-form';
import { useLoginMutation } from '../../hooks/mutations/useLoginMutation';
import Button from '../../components/Button/Button';
import Input from '../../components/Input/Input';
import Spinner from '../../components/Spinner/Spinner';
import './LoginPage.scss';

interface LoginFormData {
  password: string;
  username: string;
}

const LoginPage = () => {
  const {
    formState: { errors },
    handleSubmit,
    register
  } = useForm<LoginFormData>();

  const { isPending, mutate } = useLoginMutation();

  const handleLogin: SubmitHandler<LoginFormData> = async data => {
    const { password, username } = data;

    await mutate({ password, username });
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
        <Button className="form-submit-btn" disabled={isPending} type="submit">
          {isPending ? <Spinner /> : 'Log In'}
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
