import { ActionMeta, GroupBase, Props as SelectProps } from 'react-select';
import AsyncSelect from 'react-select/async';
import classNames from 'classnames';
import '../Dropdown.scss';

interface DropdownProps<TOption> {
  label?: string;
  loadOptions: () => Promise<TOption[]>;
  menuPortalSelector?: string;
  noOptionsMsg?: string;
  onChange: (option: TOption, e: ActionMeta<TOption>) => void;
  value?: TOption | null;
}

const AsyncDropdown = <
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
  loadOptions,
  menuPortalSelector,
  noOptionsMsg = 'No options available',
  onChange,
  placeholder = 'Search...',
  value,
  ...rest
}: SelectProps<IOption, TIsMulti, UGroup> & DropdownProps<IOption>) => {
  const portalElement = menuPortalSelector
    ? document.querySelector<HTMLElement>(menuPortalSelector)
    : undefined;

  return (
    <div className="dropdown-container">
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
      <AsyncSelect
        className={classNames('dropdown-select', {
          [className as string]: className
        })}
        classNamePrefix={classNames('dropdown-select', {
          [classNamePrefix as string]: classNamePrefix
        })}
        closeMenuOnSelect
        components={{ IndicatorSeparator: () => null, ...components }}
        defaultOptions={true}
        loadOptions={loadOptions}
        isDisabled={isDisabled}
        isOptionDisabled={option => Boolean(option.disabled)}
        id={id}
        menuPlacement="bottom"
        menuPortalTarget={portalElement}
        noOptionsMessage={() => noOptionsMsg}
        onChange={onChange}
        placeholder={placeholder}
        value={value}
        {...rest}
      />
    </div>
  );
};

export default AsyncDropdown;
