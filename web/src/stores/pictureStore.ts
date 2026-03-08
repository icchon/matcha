import { create } from 'zustand';
import type { Picture } from '@/types';
import * as picturesApi from '@/api/pictures';
import { toUserFacingMessage } from '@/lib/errorUtils';
import { MAX_PICTURES } from '@/lib/constants';

interface PictureState {
  readonly pictures: readonly Picture[];
  readonly isLoading: boolean;
  readonly error: string | null;
}

interface PictureActions {
  readonly uploadPicture: (file: File) => Promise<void>;
  readonly deletePicture: (pictureId: number) => Promise<void>;
  readonly clearError: () => void;
}

type PictureStore = PictureState & PictureActions;

const initialState: PictureState = {
  pictures: [],
  isLoading: false,
  error: null,
};

export const usePictureStore = create<PictureStore>()((set, get) => ({
  ...initialState,

  uploadPicture: async (file: File) => {
    const { pictures } = get();
    if (pictures.length >= MAX_PICTURES) {
      set({ error: 'Maximum 5 pictures allowed' });
      return;
    }
    set({ isLoading: true, error: null });
    try {
      const picture = await picturesApi.uploadPicture(file);
      set({ pictures: [...get().pictures, picture], isLoading: false });
    } catch (err) {
      const message = toUserFacingMessage(err, 'Failed to upload picture');
      set({ error: message, isLoading: false });
    }
  },

  deletePicture: async (pictureId: number) => {
    set({ isLoading: true, error: null });
    try {
      await picturesApi.deletePicture(pictureId);
      set({
        pictures: get().pictures.filter((p) => p.id !== pictureId),
        isLoading: false,
      });
    } catch (err) {
      const message = toUserFacingMessage(err, 'Failed to delete picture');
      set({ error: message, isLoading: false });
    }
  },

  clearError: () => {
    set({ error: null });
  },
}));
