import { ComponentProps } from 'react';
import classNames from 'classnames';
import './Input.scss';

interface InputProps extends ComponentProps<'input'> {
  error?: { message: string };
  label: string;
}

const Input = ({
  disabled,
  error,
  id,
  label,
  onChange,
  type = 'text',
  value,
  ...rest
}: InputProps) => (
  <label
    className={classNames('input-container', {
      error: Boolean(error)
    })}
    htmlFor={id}
  >
    {label && <span className="input-label">{label}</span>}
    <input
      className="input-field"
      disabled={disabled}
      id={id}
      onChange={onChange}
      type={type}
      value={value}
      {...rest}
    />
    {error && <span className="input-error-message">{error.message}</span>}
  </label>
);

export default Input;
