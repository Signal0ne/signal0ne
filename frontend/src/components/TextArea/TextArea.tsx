import type { ComponentProps } from 'react';
import classNames from 'classnames';
import './TextArea.scss';

interface TextAreaProps extends ComponentProps<'textarea'> {
  canResize?: boolean;
  error?: { message: string };
  label: string;
}

const TextArea = ({
  disabled,
  error,
  id,
  label,
  onChange,
  value,
  ...rest
}: TextAreaProps) => (
  <label
    className={classNames('textarea-container', {
      error: Boolean(error)
    })}
    htmlFor={id}
  >
    {label && <span className="textarea-label">{label}</span>}
    <textarea
      className={classNames('textarea-field', {
        'no-resize': !rest.canResize
      })}
      disabled={disabled}
      id={id}
      onChange={onChange}
      value={value}
      {...rest}
    />
    {error && <span className="textarea-error-message">{error.message}</span>}
  </label>
);

export default TextArea;
