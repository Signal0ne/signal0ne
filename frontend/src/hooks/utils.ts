import type { RefreshAccessTokenFn } from '../contexts/AuthProvider/AuthProvider';

interface FetchDataWithAuthProps {
  options: RequestInit;
  refreshAccessToken: RefreshAccessTokenFn;
  url: string;
}

async function retryWithRefreshedAccessToken<T>({
  options,
  refreshAccessToken,
  url
}: FetchDataWithAuthProps) {
  const accessToken = await refreshAccessToken();

  if (accessToken) {
    const retryResponse = await fetch(url, {
      ...options,
      headers: {
        ...options?.headers,
        Authorization: `Bearer ${accessToken}`
      }
    });

    if (!retryResponse.ok) {
      throw new Error('Error retrying request after refreshing token');
    }

    const data: T = await retryResponse.json();

    return data;
  } else {
    throw new Error('Token refresh failed');
  }
}

export async function fetchDataWithAuth<T>({
  options,
  refreshAccessToken,
  url
}: FetchDataWithAuthProps): Promise<T> {
  const response = await fetch(url, options);

  if (!response.ok) {
    if (response.status === 401) {
      return await retryWithRefreshedAccessToken({
        options,
        refreshAccessToken,
        url
      });
    } else {
      throw new Error('Something went wrong.');
    }
  }

  const data: T = await response.json();

  return data;
}
