import type { ComponentProps } from 'react';
import { SearchIcon } from '../Icons/Icons';
import './SearchInput.scss';

interface SearchInputProps extends ComponentProps<'input'> {}

const SearchInput = ({
  disabled,
  id,
  onChange,
  placeholder = 'Search...',
  value
}: SearchInputProps) => (
  <div className="search-input-container">
    <input
      aria-disabled={Boolean(disabled)}
      className="search-input"
      disabled={disabled}
      id={id}
      onChange={onChange}
      placeholder={placeholder}
      type="text"
      value={value}
    />
    <SearchIcon className="search-icon" height={28} width={28} />
  </div>
);

export default SearchInput;
