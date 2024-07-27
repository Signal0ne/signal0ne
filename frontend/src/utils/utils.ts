/* eslint-disable @typescript-eslint/no-explicit-any */
export const handleKeyDown =
  (callback: (...arg: any) => void, disabled?: boolean) =>
  (e: React.KeyboardEvent) => {
    e.key === ' ' && e.preventDefault();

    if (['Enter', ' '].includes(e.key) && !disabled) callback(e);
  };
