import { useEffect, type FC } from 'react';
import { toast } from 'sonner';
import { useProfileStore } from '@/stores/profileStore';
import { usePictureStore } from '@/stores/pictureStore';
import { useTagStore } from '@/stores/tagStore';
import { ProfileForm } from '@/features/profile/components/ProfileForm';
import { PhotoUploader } from '@/features/profile/components/PhotoUploader';
import { TagManager } from '@/features/profile/components/TagManager';
import { Spinner } from '@/components/ui/Spinner';
import type { ProfileFormData } from '@/lib/validators';

const ProfilePage: FC = () => {
  const profile = useProfileStore((s) => s.profile);
  const isLoading = useProfileStore((s) => s.isLoading);
  const error = useProfileStore((s) => s.error);
  const fetchProfile = useProfileStore((s) => s.fetchProfile);
  const createProfile = useProfileStore((s) => s.createProfile);
  const updateProfile = useProfileStore((s) => s.updateProfile);

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
    fetchProfile();
    fetchTags();
  }, [fetchProfile, fetchTags]);

  const isNewProfile = !isLoading && profile === null && error === null;

  const handleSubmit = async (data: ProfileFormData) => {
    // Read current profile at call time to avoid stale closure if user double-submits
    const currentProfile = useProfileStore.getState().profile;
    const isCreate = currentProfile === null;

    const params = {
      firstName: data.firstName,
      lastName: data.lastName,
      username: data.username,
      gender: data.gender,
      sexualPreference: data.sexualPreference,
      birthday: data.birthday,
      biography: data.biography,
      occupation: data.occupation || undefined,
    };

    const success = isCreate
      ? await createProfile(params)
      : await updateProfile(params);

    if (success) {
      toast.success(isCreate ? 'Profile created!' : 'Profile updated!');
    }
  };

  if (isLoading && !profile && !error) {
    return (
      <div className="flex items-center justify-center p-12">
        <Spinner size="lg" />
      </div>
    );
  }

  if (error) {
    const errors = [error, pictureError, tagError].filter(Boolean);
    return (
      <div className="mx-auto max-w-2xl space-y-8 p-6">
        <div role="alert" className="rounded-md bg-red-50 p-4 text-sm text-red-700">
          {errors.length === 1 ? errors[0] : (
            <ul className="list-disc pl-4">
              {errors.map((e, i) => <li key={i}>{e}</li>)}
            </ul>
          )}
        </div>
      </div>
    );
  }

  const nonProfileErrors = [pictureError, tagError].filter(Boolean);

  return (
    <div className="mx-auto max-w-2xl space-y-8 p-6">
      <h1 className="text-2xl font-bold text-gray-900">
        {isNewProfile ? 'Create Profile' : 'Edit Profile'}
      </h1>

      {nonProfileErrors.length > 0 ? (
        <div role="alert" className="rounded-md bg-red-50 p-4 text-sm text-red-700">
          {nonProfileErrors.length === 1 ? nonProfileErrors[0] : (
            <ul className="list-disc pl-4">
              {nonProfileErrors.map((e, i) => <li key={i}>{e}</li>)}
            </ul>
          )}
        </div>
      ) : null}

      <ProfileForm
        onSubmit={handleSubmit}
        isLoading={isLoading}
        initialValues={profile ? {
          firstName: profile.firstName ?? undefined,
          lastName: profile.lastName ?? undefined,
          username: profile.username ?? undefined,
          gender: profile.gender ?? undefined,
          sexualPreference: profile.sexualPreference ?? undefined,
          birthday: profile.birthday ?? undefined,
          biography: profile.biography ?? undefined,
          occupation: profile.occupation ?? undefined,
        } : undefined}
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

export { ProfilePage };
