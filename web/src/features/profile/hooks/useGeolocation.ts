import { useState, useCallback } from 'react';

interface UseGeolocationReturn {
  readonly latitude: number | null;
  readonly longitude: number | null;
  readonly error: string | null;
  readonly isLoading: boolean;
  readonly requestLocation: () => void;
  readonly setManualLocation: (lat: number, lng: number) => void;
}

// TODO(FE-GEOLOCATION): Wire useGeolocation to the edit profile page so users can update their location.
export function useGeolocation(): UseGeolocationReturn {
  const [latitude, setLatitude] = useState<number | null>(null);
  const [longitude, setLongitude] = useState<number | null>(null);
  const [error, setError] = useState<string | null>(null);
  const [isLoading, setIsLoading] = useState(false);

  const requestLocation = useCallback(() => {
    if (!navigator.geolocation) {
      setError('Geolocation is not supported by this browser');
      return;
    }
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
    if (!Number.isFinite(lat) || !Number.isFinite(lng) || lat < -90 || lat > 90 || lng < -180 || lng > 180) {
      setError('Invalid coordinates: latitude must be -90..90, longitude must be -180..180');
      return;
    }
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
