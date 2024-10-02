import classNames from 'classnames';
import Select, {
  ActionMeta,
  GroupBase,
  Props as SelectProps
} from 'react-select';
import './Dropdown.scss';

interface DropdownProps<TOption> {
  label?: string;
  menuPortalSelector?: string;
  noOptionsMsg?: string;
  onChange: (option: TOption, e: ActionMeta<TOption>) => void;
  options?: TOption[];
  value?: TOption | null;
}

const Dropdown = <
  IOption extends { disabled?: boolean },
  TIsMulti extends boolean = false,
  UGroup extends GroupBase<IOption> = GroupBase<IOption>
>({
  className,
  classNamePrefix,
  components,
  id,
  isDisabled,
  label,
  menuPortalSelector,
  noOptionsMsg = 'No options available',
  onChange,
  options,
  placeholder = 'Search...',
  value,
  ...rest
}: SelectProps<IOption, TIsMulti, UGroup> & DropdownProps<IOption>) => {
  const portalElement = menuPortalSelector
    ? document.querySelector<HTMLElement>(menuPortalSelector)
    : undefined;

  return (
    <div
      className={classNames('dropdown-container', {
        'is-disabled': isDisabled
      })}
    >
      {label && (
        <label
          className={classNames('dropdown-label', {
            'is-disabled': isDisabled
          })}
          htmlFor={id}
        >
          {label}
        </label>
      )}
      <Select
        className={classNames('dropdown-select', {
          [className as string]: className
        })}
        classNamePrefix={classNames('dropdown-select', {
          [classNamePrefix as string]: classNamePrefix
        })}
        closeMenuOnSelect
        components={{ IndicatorSeparator: () => null, ...components }}
        id={id}
        isDisabled={isDisabled}
        isOptionDisabled={option => Boolean(option.disabled)}
        menuPlacement="bottom"
        menuPortalTarget={portalElement}
        noOptionsMessage={() => noOptionsMsg}
        onChange={onChange}
        options={options}
        placeholder={placeholder}
        value={value}
        {...rest}
      />
    </div>
  );
};

export default Dropdown;
