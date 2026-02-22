import { useEffect, type FC } from 'react';
import { useProfileStore } from '@/stores/profileStore';
import { usePictureStore } from '@/stores/pictureStore';
import { useTagStore } from '@/stores/tagStore';
import { ProfileForm } from '@/features/profile/components/ProfileForm';
import { PhotoUploader } from '@/features/profile/components/PhotoUploader';
import { TagManager } from '@/features/profile/components/TagManager';
import type { ProfileFormData } from '@/lib/validators';

const ProfileCreatePage: FC = () => {
  const profile = useProfileStore((s) => s.profile);
  const isProfileLoading = useProfileStore((s) => s.isLoading);
  const profileError = useProfileStore((s) => s.error);
  const createProfile = useProfileStore((s) => s.createProfile);

  const pictures = usePictureStore((s) => s.pictures);
  const isPicturesLoading = usePictureStore((s) => s.isLoading);
  const pictureError = usePictureStore((s) => s.error);
  const uploadPicture = usePictureStore((s) => s.uploadPicture);
  const deletePicture = usePictureStore((s) => s.deletePicture);

  const tags = useTagStore((s) => s.tags);
  const allTags = useTagStore((s) => s.allTags);
  const isTagsLoading = useTagStore((s) => s.isLoading);
  const tagError = useTagStore((s) => s.error);
  const fetchTags = useTagStore((s) => s.fetchTags);
  const addTag = useTagStore((s) => s.addTag);
  const removeTag = useTagStore((s) => s.removeTag);

  useEffect(() => {
    fetchTags();
  }, [fetchTags]);

  const handleSubmit = (data: ProfileFormData) => {
    createProfile({
      firstName: data.firstName,
      lastName: data.lastName,
      username: data.username,
      gender: data.gender,
      sexualPreference: data.sexualPreference,
      birthday: data.birthday,
      biography: data.biography,
      occupation: data.occupation || undefined,
    });
  };

  const errors = [profileError, pictureError, tagError].filter(Boolean);

  return (
    <div className="mx-auto max-w-2xl space-y-8 p-6">
      <h1 className="text-2xl font-bold text-gray-900">Create Profile</h1>

      {errors.length > 0 ? (
        <div role="alert" className="rounded-md bg-red-50 p-4 text-sm text-red-700">
          {errors.length === 1 ? errors[0] : (
            <ul className="list-disc pl-4">
              {errors.map((e) => <li key={e}>{e}</li>)}
            </ul>
          )}
        </div>
      ) : null}

      <ProfileForm
        onSubmit={handleSubmit}
        isLoading={isProfileLoading}
        initialValues={profile ?? undefined}
      />

      <PhotoUploader
        pictures={pictures}
        onUpload={uploadPicture}
        onDelete={deletePicture}
        isLoading={isPicturesLoading}
      />

      <TagManager
        tags={tags}
        allTags={allTags}
        onAdd={addTag}
        onRemove={removeTag}
        isLoading={isTagsLoading}
      />
    </div>
  );
};

export { ProfileCreatePage };
