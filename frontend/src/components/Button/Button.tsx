import { ComponentProps } from 'react';
import classNames from 'classnames';
import './Button.scss';

interface ButtonProps extends ComponentProps<'button'> {
  className?: string;
}

const Button = ({
  children,
  className,
  disabled,
  id,
  onClick,
  type = 'button'
}: ButtonProps) => (
  <button
    aria-disabled={disabled}
    className={classNames('button', {
      [className as string]: className
    })}
    disabled={disabled}
    id={id}
    onClick={onClick}
    type={type}
  >
    {children}
  </button>
);

export default Button;
