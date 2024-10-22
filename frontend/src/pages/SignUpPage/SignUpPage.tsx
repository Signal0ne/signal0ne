import { Link } from 'react-router-dom';
import { SubmitHandler, useForm } from 'react-hook-form';
import { useRegisterMutation } from '../../hooks/mutations/useRegisterMutation';
import Button from '../../components/Button/Button';
import Input from '../../components/Input/Input';
import Spinner from '../../components/Spinner/Spinner';
import './SignUpPage.scss';

interface SignUpFormData {
  confirmPassword: string;
  password: string;
  username: string;
}

const SignUpPage = () => {
  const {
    formState: { errors },
    handleSubmit,
    register,
    watch
  } = useForm<SignUpFormData>();

  const password = watch('password');

  const { isPending, mutate } = useRegisterMutation();

  const handleSignUp: SubmitHandler<SignUpFormData> = async data => {
    const { password, username } = data;

    await mutate({ password, username });
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
                    'Password must contain at least one uppercase, lowercase letter, one number and one special character',
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
            {errors.confirmPassword && (
              <span className="error-msg">
                {errors.confirmPassword?.message}
              </span>
            )}
          </div>
        </div>
        <Button className="form-submit-btn" disabled={isPending} type="submit">
          {isPending ? <Spinner /> : 'Register'}
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
