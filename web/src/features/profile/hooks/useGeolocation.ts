import { useState, useCallback } from 'react';

interface UseGeolocationReturn {
  readonly latitude: number | null;
  readonly longitude: number | null;
  readonly error: string | null;
  readonly isLoading: boolean;
  readonly requestLocation: () => void;
  readonly setManualLocation: (lat: number, lng: number) => void;
}

export function useGeolocation(): UseGeolocationReturn {
  const [latitude, setLatitude] = useState<number | null>(null);
  const [longitude, setLongitude] = useState<number | null>(null);
  const [error, setError] = useState<string | null>(null);
  const [isLoading, setIsLoading] = useState(false);

  const requestLocation = useCallback(() => {
    setIsLoading(true);
    setError(null);
    navigator.geolocation.getCurrentPosition(
      (position) => {
        setLatitude(position.coords.latitude);
        setLongitude(position.coords.longitude);
        setIsLoading(false);
      },
      (err) => {
        setError(err.message);
        setIsLoading(false);
      },
    );
  }, []);

  const setManualLocation = useCallback((lat: number, lng: number) => {
    setLatitude(lat);
    setLongitude(lng);
    setError(null);
  }, []);

  return {
    latitude,
    longitude,
    error,
    isLoading,
    requestLocation,
    setManualLocation,
  };
}
