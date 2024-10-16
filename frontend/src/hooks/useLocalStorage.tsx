import { useState } from 'react';

export const useLocalStorage = (keyName: string, defaultValue: unknown) => {
  const [storedValue, setStoredValue] = useState(() => {
    try {
      const value = localStorage.getItem(keyName);

      if (value) {
        return JSON.parse(value);
      } else {
        localStorage.setItem(keyName, JSON.stringify(defaultValue));
        return defaultValue;
      }
    } catch (err) {
      return defaultValue;
    }
  });

  const removeValue = () => {
    try {
      localStorage.removeItem(keyName);
    } catch (err) {
      console.error(err);
    }
  };

  const setValue = (newValue: unknown) => {
    try {
      localStorage.setItem(keyName, JSON.stringify(newValue));
    } catch (err) {
      console.error(err);
    }

    setStoredValue(newValue);
  };

  return { removeValue, setValue, storedValue };
};
