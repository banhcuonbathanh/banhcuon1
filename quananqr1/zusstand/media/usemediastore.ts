import { create } from 'zustand';
import { persist } from 'zustand/middleware';

import { UploadImageResType } from '@/schemaValidations/media.schema';
import { useApiStore } from '../api/api-controller';
interface MediaStore {
  isUploading: boolean;
  error: string | null;
  uploadMedia: (file: File) => Promise<string>;
}

export const useMediaStore = create<MediaStore>()(
  persist(
    (set, get) => ({
      isUploading: false,
      error: null,
      uploadMedia: async (file: File) => {
        const { http } = useApiStore.getState();
        set({ isUploading: true, error: null });
        try {
          const formData = new FormData();
          formData.append('file', file);
          const response = await http.post<UploadImageResType>('/media/upload', formData);
          set({ isUploading: false });
          return response.data.data; // Assuming the response contains the image URL in data.data
        } catch (error) {
          set({ isUploading: false, error: "Failed to upload media" });
          throw error;
        }
      },
    }),
    {
      name: "media-storage",
      skipHydration: true
    }
  )
);

// Custom hook for media upload mutation
export const useUploadMediaMutation = () => {
  const { uploadMedia, isUploading, error } = useMediaStore();
  return {
    mutateAsync: uploadMedia,
    isPending: isUploading,
    error
  };
};