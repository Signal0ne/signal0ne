import { ComponentProps, forwardRef } from 'react';
import classNames from 'classnames';
import './Input.scss';

interface InputProps extends ComponentProps<'input'> {
  error?: { message: string };
  label: string;
}

const Input = forwardRef<HTMLInputElement, InputProps>(
  (
    { disabled, error, id, label, onChange, type = 'text', value, ...rest },
    ref
  ) => (
    <label
      className={classNames('input-container', {
        error: Boolean(error)
      })}
      htmlFor={id}
    >
      {label && <span className="input-label">{label}</span>}
      <input
        aria-disabled={disabled}
        className="input-field"
        disabled={disabled}
        id={id}
        onChange={onChange}
        ref={ref}
        type={type}
        value={value}
        {...rest}
      />
      {error && <span className="input-error-message">{error.message}</span>}
    </label>
  )
);

export default Input;
