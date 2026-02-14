import { describe, it, expect, vi, beforeEach } from 'vitest';
import { renderHook, act } from '@testing-library/react';
import { useGeolocation } from '@/features/profile/hooks/useGeolocation';

const mockGetCurrentPosition = vi.fn();

beforeEach(() => {
  vi.resetAllMocks();
  Object.defineProperty(globalThis.navigator, 'geolocation', {
    value: { getCurrentPosition: mockGetCurrentPosition },
    writable: true,
    configurable: true,
  });
});

describe('useGeolocation', () => {
  it('has null coordinates initially', () => {
    const { result } = renderHook(() => useGeolocation());

    expect(
      result.current.latitude,
      'Initial latitude should be null before any geolocation request.',
    ).toBeNull();
    expect(
      result.current.longitude,
      'Initial longitude should be null before any geolocation request.',
    ).toBeNull();
    expect(
      result.current.error,
      'Initial error should be null.',
    ).toBeNull();
    expect(
      result.current.isLoading,
      'Initial isLoading should be false.',
    ).toBe(false);
  });

  it('requestLocation sets coordinates from Geolocation API', async () => {
    mockGetCurrentPosition.mockImplementation((success: PositionCallback) => {
      success({
        coords: { latitude: 35.6762, longitude: 139.6503 },
      } as GeolocationPosition);
    });

    const { result } = renderHook(() => useGeolocation());

    act(() => {
      result.current.requestLocation();
    });

    expect(
      result.current.latitude,
      'latitude should be set from Geolocation API response.',
    ).toBe(35.6762);
    expect(
      result.current.longitude,
      'longitude should be set from Geolocation API response.',
    ).toBe(139.6503);
    expect(
      result.current.isLoading,
      'isLoading should be false after successful geolocation.',
    ).toBe(false);
  });

  it('sets error when Geolocation API fails', () => {
    mockGetCurrentPosition.mockImplementation(
      (_success: PositionCallback, error: PositionErrorCallback) => {
        error({
          code: 1,
          message: 'User denied Geolocation',
          PERMISSION_DENIED: 1,
          POSITION_UNAVAILABLE: 2,
          TIMEOUT: 3,
        } as GeolocationPositionError);
      },
    );

    const { result } = renderHook(() => useGeolocation());

    act(() => {
      result.current.requestLocation();
    });

    expect(
      result.current.error,
      'error should be set when Geolocation API is denied. Check error callback handling.',
    ).toBe('User denied Geolocation');
    expect(
      result.current.latitude,
      'latitude should remain null on error.',
    ).toBeNull();
  });

  it('sets error when geolocation is not supported', () => {
    Object.defineProperty(globalThis.navigator, 'geolocation', {
      value: undefined,
      writable: true,
      configurable: true,
    });

    const { result } = renderHook(() => useGeolocation());

    act(() => {
      result.current.requestLocation();
    });

    expect(
      result.current.error,
      'error should indicate geolocation is unsupported when navigator.geolocation is undefined.',
    ).toBe('Geolocation is not supported by this browser');
    expect(
      result.current.isLoading,
      'isLoading should remain false when geolocation is unsupported.',
    ).toBe(false);
  });

  it('setManualLocation sets coordinates directly', () => {
    const { result } = renderHook(() => useGeolocation());

    act(() => {
      result.current.setManualLocation(48.8566, 2.3522);
    });

    expect(
      result.current.latitude,
      'setManualLocation should set latitude directly.',
    ).toBe(48.8566);
    expect(
      result.current.longitude,
      'setManualLocation should set longitude directly.',
    ).toBe(2.3522);
    expect(
      result.current.error,
      'error should be cleared after manual location set.',
    ).toBeNull();
  });

  it('setManualLocation clears previous error', () => {
    mockGetCurrentPosition.mockImplementation(
      (_success: PositionCallback, error: PositionErrorCallback) => {
        error({
          code: 1,
          message: 'Denied',
          PERMISSION_DENIED: 1,
          POSITION_UNAVAILABLE: 2,
          TIMEOUT: 3,
        } as GeolocationPositionError);
      },
    );

    const { result } = renderHook(() => useGeolocation());

    act(() => {
      result.current.requestLocation();
    });
    expect(result.current.error).not.toBeNull();

    act(() => {
      result.current.setManualLocation(40.7128, -74.006);
    });

    expect(
      result.current.error,
      'setManualLocation should clear any previous geolocation error.',
    ).toBeNull();
    expect(result.current.latitude).toBe(40.7128);
  });
});
