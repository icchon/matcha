import { useEffect, type FC } from 'react';
import { useProfileStore } from '@/stores/profileStore';
import { ProfileForm } from '@/features/profile/components/ProfileForm';
import { PhotoUploader } from '@/features/profile/components/PhotoUploader';
import { TagManager } from '@/features/profile/components/TagManager';
import type { ProfileFormData } from '@/lib/validators';

const ProfileCreatePage: FC = () => {
  const profile = useProfileStore((s) => s.profile);
  const pictures = useProfileStore((s) => s.pictures);
  const tags = useProfileStore((s) => s.tags);
  const allTags = useProfileStore((s) => s.allTags);
  const isLoading = useProfileStore((s) => s.isLoading);
  const error = useProfileStore((s) => s.error);
  const saveProfile = useProfileStore((s) => s.saveProfile);
  const fetchTags = useProfileStore((s) => s.fetchTags);
  const uploadPicture = useProfileStore((s) => s.uploadPicture);
  const deletePicture = useProfileStore((s) => s.deletePicture);
  const addTag = useProfileStore((s) => s.addTag);
  const removeTag = useProfileStore((s) => s.removeTag);

  useEffect(() => {
    fetchTags();
  }, [fetchTags]);

  const handleSubmit = (data: ProfileFormData) => {
    saveProfile({
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

  return (
    <div className="mx-auto max-w-2xl space-y-8 p-6">
      <h1 className="text-2xl font-bold text-gray-900">Create Profile</h1>

      {error ? (
        <div role="alert" className="rounded-md bg-red-50 p-4 text-sm text-red-700">
          {error}
        </div>
      ) : null}

      <ProfileForm
        onSubmit={handleSubmit}
        isLoading={isLoading}
        initialValues={profile ?? undefined}
      />

      <PhotoUploader
        pictures={[...pictures]}
        onUpload={uploadPicture}
        onDelete={deletePicture}
        isLoading={isLoading}
      />

      <TagManager
        tags={[...tags]}
        allTags={[...allTags]}
        onAdd={addTag}
        onRemove={removeTag}
        isLoading={isLoading}
      />
    </div>
  );
};

export { ProfileCreatePage };
