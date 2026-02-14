import { describe, it, expect, vi, beforeEach } from 'vitest';
import { render, screen, waitFor } from '@testing-library/react';
import userEvent from '@testing-library/user-event';
import { PhotoUploader } from '@/features/profile/components/PhotoUploader';
import type { Picture } from '@/types';

const mockOnUpload = vi.fn();
const mockOnDelete = vi.fn();

const samplePictures: Picture[] = [
  { id: 1, userId: 'u1', url: '/images/1.jpg', isProfilePic: true, createdAt: '2025-01-01T00:00:00Z' },
  { id: 2, userId: 'u1', url: '/images/2.jpg', isProfilePic: false, createdAt: '2025-01-02T00:00:00Z' },
];

beforeEach(() => {
  vi.resetAllMocks();
});

describe('PhotoUploader', () => {
  it('renders existing pictures', () => {
    render(
      <PhotoUploader
        pictures={samplePictures}
        onUpload={mockOnUpload}
        onDelete={mockOnDelete}
        isLoading={false}
      />,
    );

    const images = screen.getAllByRole('img');
    expect(
      images,
      'PhotoUploader should render an img element for each picture.',
    ).toHaveLength(2);
  });

  it('renders upload area', () => {
    render(
      <PhotoUploader
        pictures={[]}
        onUpload={mockOnUpload}
        onDelete={mockOnDelete}
        isLoading={false}
      />,
    );

    expect(
      screen.getByText(/upload/i),
      'PhotoUploader should display an upload area.',
    ).toBeInTheDocument();
  });

  it('calls onUpload when a file is selected', async () => {
    const user = userEvent.setup();
    render(
      <PhotoUploader
        pictures={[]}
        onUpload={mockOnUpload}
        onDelete={mockOnDelete}
        isLoading={false}
      />,
    );

    const file = new File(['image-data'], 'photo.jpg', { type: 'image/jpeg' });
    const input = screen.getByTestId('photo-file-input') as HTMLInputElement;
    await user.upload(input, file);

    await waitFor(() => {
      expect(
        mockOnUpload,
        'onUpload should be called with the selected file.',
      ).toHaveBeenCalledWith(file);
    });
  });

  it('calls onDelete when delete button is clicked', async () => {
    const user = userEvent.setup();
    render(
      <PhotoUploader
        pictures={samplePictures}
        onUpload={mockOnUpload}
        onDelete={mockOnDelete}
        isLoading={false}
      />,
    );

    const deleteButtons = screen.getAllByRole('button', { name: /delete/i });
    await user.click(deleteButtons[0]);

    expect(
      mockOnDelete,
      'onDelete should be called with the picture ID when delete button is clicked.',
    ).toHaveBeenCalledWith(1);
  });

  it('hides upload area when at max 5 pictures', () => {
    const fivePictures: Picture[] = Array.from({ length: 5 }, (_, i) => ({
      id: i + 1,
      userId: 'u1',
      url: `/images/${i + 1}.jpg`,
      isProfilePic: i === 0,
      createdAt: '2025-01-01T00:00:00Z',
    }));
    render(
      <PhotoUploader
        pictures={fivePictures}
        onUpload={mockOnUpload}
        onDelete={mockOnDelete}
        isLoading={false}
      />,
    );

    expect(
      screen.queryByTestId('photo-file-input'),
      'Upload input should be hidden when at max 5 pictures.',
    ).not.toBeInTheDocument();
  });

  it('shows picture count', () => {
    render(
      <PhotoUploader
        pictures={samplePictures}
        onUpload={mockOnUpload}
        onDelete={mockOnDelete}
        isLoading={false}
      />,
    );

    expect(
      screen.getByText(/2.*5/),
      'Should display current count out of max (e.g., 2/5).',
    ).toBeInTheDocument();
  });
});
