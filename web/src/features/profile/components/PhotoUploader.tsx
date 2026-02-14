import { useCallback, useRef, type FC, type ChangeEvent, type DragEvent } from 'react';
import { useState } from 'react';
import { Button } from '@/components/ui/Button';
import type { Picture } from '@/types';

const MAX_PICTURES = 5;

interface PhotoUploaderProps {
  readonly pictures: readonly Picture[];
  readonly onUpload: (file: File) => void;
  readonly onDelete: (pictureId: number) => void;
  readonly isLoading: boolean;
}

const PhotoUploader: FC<PhotoUploaderProps> = ({ pictures, onUpload, onDelete, isLoading }) => {
  const [isDragging, setIsDragging] = useState(false);
  const fileInputRef = useRef<HTMLInputElement>(null);
  const canUpload = pictures.length < MAX_PICTURES;

  const handleFileChange = useCallback(
    (e: ChangeEvent<HTMLInputElement>) => {
      const file = e.target.files?.[0];
      if (file) {
        onUpload(file);
      }
      if (fileInputRef.current) {
        fileInputRef.current.value = '';
      }
    },
    [onUpload],
  );

  const handleDrop = useCallback(
    (e: DragEvent<HTMLDivElement>) => {
      e.preventDefault();
      setIsDragging(false);
      const file = e.dataTransfer.files[0];
      if (file) {
        onUpload(file);
      }
    },
    [onUpload],
  );

  const handleDragOver = useCallback((e: DragEvent<HTMLDivElement>) => {
    e.preventDefault();
    setIsDragging(true);
  }, []);

  const handleDragLeave = useCallback((e: DragEvent<HTMLDivElement>) => {
    e.preventDefault();
    setIsDragging(false);
  }, []);

  return (
    <div className="space-y-4">
      <div className="flex items-center justify-between">
        <h3 className="text-lg font-medium">Photos</h3>
        <span className="text-sm text-gray-500">{pictures.length}/{MAX_PICTURES}</span>
      </div>

      <div className="grid grid-cols-2 gap-4 sm:grid-cols-3">
        {pictures.map((picture) => (
          <div key={picture.id} className="group relative">
            <img
              src={picture.url}
              alt={`Photo ${picture.id}`}
              className="h-32 w-full rounded-lg object-cover"
            />
            <Button
              variant="danger"
              size="sm"
              aria-label={`Delete photo ${picture.id}`}
              className="absolute right-1 top-1 opacity-0 group-hover:opacity-100"
              onClick={() => onDelete(picture.id)}
              disabled={isLoading}
            >
              Delete
            </Button>
          </div>
        ))}
      </div>

      {canUpload ? (
        <div
          onDrop={handleDrop}
          onDragOver={handleDragOver}
          onDragLeave={handleDragLeave}
          className={`flex cursor-pointer flex-col items-center justify-center rounded-lg border-2 border-dashed p-6 transition-colors ${
            isDragging ? 'border-blue-500 bg-blue-50' : 'border-gray-300 hover:border-gray-400'
          }`}
          onClick={() => fileInputRef.current?.click()}
        >
          <p className="text-sm text-gray-600">
            Drag & drop or click to upload
          </p>
          <input
            ref={fileInputRef}
            type="file"
            accept="image/*"
            data-testid="photo-file-input"
            className="hidden"
            onChange={handleFileChange}
            disabled={isLoading}
          />
        </div>
      ) : null}
    </div>
  );
};

export { PhotoUploader };
