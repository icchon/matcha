import type { FC } from 'react';
import { useForm } from 'react-hook-form';
import { zodResolver } from '@hookform/resolvers/zod';
import { profileSchema, type ProfileFormData } from '@/lib/validators';
import { Button } from '@/components/ui/Button';
import { Input } from '@/components/ui/Input';
import type { Gender, SexualPreference } from '@/types';

interface ProfileFormProps {
  readonly onSubmit: (data: ProfileFormData) => void;
  readonly isLoading: boolean;
  readonly initialValues?: {
    readonly firstName?: string;
    readonly lastName?: string;
    readonly username?: string;
    readonly gender?: Gender;
    readonly sexualPreference?: SexualPreference;
    readonly birthday?: string;
    readonly biography?: string;
    readonly occupation?: string;
  };
}

const ProfileForm: FC<ProfileFormProps> = ({ onSubmit, isLoading, initialValues }) => {
  const {
    register,
    handleSubmit,
    formState: { errors },
  } = useForm<ProfileFormData>({
    resolver: zodResolver(profileSchema),
    defaultValues: {
      firstName: initialValues?.firstName ?? '',
      lastName: initialValues?.lastName ?? '',
      username: initialValues?.username ?? '',
      gender: initialValues?.gender,
      sexualPreference: initialValues?.sexualPreference,
      birthday: initialValues?.birthday ?? '',
      biography: initialValues?.biography ?? '',
      occupation: initialValues?.occupation ?? '',
    },
  });

  return (
    <form onSubmit={handleSubmit(onSubmit)} noValidate className="flex flex-col gap-4">
      <Input
        label="First Name"
        error={errors.firstName?.message}
        {...register('firstName')}
      />
      <Input
        label="Last Name"
        error={errors.lastName?.message}
        {...register('lastName')}
      />
      <Input
        label="Username"
        error={errors.username?.message}
        {...register('username')}
      />

      <div className="flex flex-col gap-1">
        <label htmlFor="gender" className="text-sm font-medium text-gray-700">
          Gender
        </label>
        <select
          id="gender"
          className="rounded-md border border-gray-300 px-3 py-2 text-sm shadow-sm focus:outline-none focus:ring-2 focus:ring-blue-500"
          {...register('gender')}
        >
          <option value="">Select gender</option>
          <option value="male">Male</option>
          <option value="female">Female</option>
          <option value="other">Other</option>
        </select>
        {errors.gender?.message ? (
          <p role="alert" className="text-sm text-red-600">{errors.gender.message}</p>
        ) : null}
      </div>

      <div className="flex flex-col gap-1">
        <label htmlFor="sexualPreference" className="text-sm font-medium text-gray-700">
          Sexual Preference
        </label>
        <select
          id="sexualPreference"
          className="rounded-md border border-gray-300 px-3 py-2 text-sm shadow-sm focus:outline-none focus:ring-2 focus:ring-blue-500"
          {...register('sexualPreference')}
        >
          <option value="">Select preference</option>
          <option value="heterosexual">Heterosexual</option>
          <option value="homosexual">Homosexual</option>
          <option value="bisexual">Bisexual</option>
        </select>
        {errors.sexualPreference?.message ? (
          <p role="alert" className="text-sm text-red-600">{errors.sexualPreference.message}</p>
        ) : null}
      </div>

      <Input
        label="Birthday"
        type="date"
        error={errors.birthday?.message}
        {...register('birthday')}
      />

      <div className="flex flex-col gap-1">
        <label htmlFor="biography" className="text-sm font-medium text-gray-700">
          Biography
        </label>
        <textarea
          id="biography"
          rows={4}
          className="rounded-md border border-gray-300 px-3 py-2 text-sm shadow-sm focus:outline-none focus:ring-2 focus:ring-blue-500"
          {...register('biography')}
        />
        {errors.biography?.message ? (
          <p role="alert" className="text-sm text-red-600">{errors.biography.message}</p>
        ) : null}
      </div>

      <Input
        label="Occupation"
        error={errors.occupation?.message}
        {...register('occupation')}
      />

      <Button type="submit" loading={isLoading}>
        Save Profile
      </Button>
    </form>
  );
};

export { ProfileForm };
