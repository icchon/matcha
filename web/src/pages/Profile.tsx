import React, { useState, useEffect, ChangeEvent, FormEvent } from 'react';
import { getMyProfile, updateMyProfile, getMyPictures, uploadPicture, setProfilePic, deletePicture, UserProfile, Picture } from '../services/profileService';

const Profile: React.FC = () => {
  const [profile, setProfile] = useState<UserProfile | null>(null);
  const [pictures, setPictures] = useState<Picture[]>([]);
  const [selectedFile, setSelectedFile] = useState<File | null>(null);
  const [loading, setLoading] = useState(true);

  const fetchProfile = async () => {
    try {
      const response = await getMyProfile();
      setProfile(response.data);
    } catch (error) {
      console.error('Error fetching profile', error);
    }
  };

  const fetchPictures = async () => {
    try {
        const response = await getMyPictures();
        setPictures(response.data);
    } catch (error) {
        console.error('Error fetching pictures', error);
    }
  }

  useEffect(() => {
    const loadData = async () => {
        setLoading(true);
        await Promise.all([fetchProfile(), fetchPictures()]);
        setLoading(false);
    }
    loadData();
  }, []);

  const handleChange = (e: ChangeEvent<HTMLInputElement | HTMLTextAreaElement>) => {
    if (profile) {
      setProfile({
        ...profile,
        [e.target.name]: e.target.value,
      });
    }
  };

  const handleSubmit = async (e: FormEvent) => {
    e.preventDefault();
    if (profile) {
      try {
        await updateMyProfile(profile);
        alert('Profile updated successfully');
      } catch (error) {
        alert('Failed to update profile');
      }
    }
  };

  const handleFileChange = (e: ChangeEvent<HTMLInputElement>) => {
    if (e.target.files) {
        setSelectedFile(e.target.files[0]);
    }
  }

  const handlePictureUpload = async (e: FormEvent) => {
    e.preventDefault();
    if (!selectedFile) {
        alert('Please select a file to upload');
        return;
    }
    const formData = new FormData();
    formData.append('image', selectedFile);
    try {
        await uploadPicture(formData);
        alert('Picture uploaded successfully');
        fetchPictures(); // Refresh pictures
        setSelectedFile(null);
    } catch (error) {
        alert('Failed to upload picture');
    }
  }

  const handleSetProfilePic = async (pictureId: number) => {
    try {
        await setProfilePic(pictureId);
        alert('Profile picture updated');
        fetchPictures(); // Refresh pictures
    } catch (error) {
        alert('Failed to set profile picture');
    }
  }

  const handleDeletePicture = async (pictureId: number) => {
    try {
        await deletePicture(pictureId);
        alert('Picture deleted');
        fetchPictures(); // Refresh pictures
    } catch (error) {
        alert('Failed to delete picture');
    }
  }

  if (loading) {
    return <div>Loading...</div>;
  }

  if (!profile) {
    return <div>Failed to load profile.</div>;
  }

  return (
    <div>
      <h1>Profile</h1>
      <div>
        <h2>My Pictures</h2>
        <div style={{ display: 'flex', flexWrap: 'wrap', gap: '10px' }}>
            {pictures.map(pic => (
                <div key={pic.id} style={{ border: pic.is_profile_pic ? '2px solid blue' : 'none', padding: '5px' }}>
                    <img src={pic.url} alt="User" width="150" />
                    <div>
                        {!pic.is_profile_pic && <button onClick={() => handleSetProfilePic(pic.id)}>Set as Profile</button>}
                        <button onClick={() => handleDeletePicture(pic.id)}>Delete</button>
                    </div>
                </div>
            ))}
        </div>
        <h3>Upload New Picture</h3>
        <form onSubmit={handlePictureUpload}>
            <input type="file" onChange={handleFileChange} />
            <button type="submit">Upload</button>
        </form>
      </div>
      <hr />
      <form onSubmit={handleSubmit}>
        <h2>Edit Profile Details</h2>
        <div>
          <label>First Name:</label>
          <input
            type="text"
            name="first_name"
            value={profile.first_name || ''}
            onChange={handleChange}
          />
        </div>
        <div>
          <label>Last Name:</label>
          <input
            type="text"
            name="last_name"
            value={profile.last_name || ''}
            onChange={handleChange}
          />
        </div>
        <div>
          <label>Username:</label>
          <input
            type="text"
            name="username"
            value={profile.username || ''}
            onChange={handleChange}
          />
        </div>
        <div>
          <label>Gender:</label>
          <input
            type="text"
            name="gender"
            value={profile.gender || ''}
            onChange={handleChange}
          />
        </div>
        <div>
          <label>Sexual Preference:</label>
          <input
            type="text"
            name="sexual_preference"
            value={profile.sexual_preference || ''}
            onChange={handleChange}
          />
        </div>
        <div>
          <label>Birthday:</label>
          <input
            type="text"
            name="birthday"
            value={profile.birthday || ''}
            onChange={handleChange}
          />
        </div>
        <div>
          <label>Occupation:</label>
          <input
            type="text"
            name="occupation"
            value={profile.occupation || ''}
            onChange={handleChange}
          />
        </div>
        <div>
          <label>Biography:</label>
          <textarea
            name="biography"
            value={profile.biography || ''}
            onChange={handleChange}
          />
        </div>
        <div>
          <label>Location:</label>
          <input
            type="text"
            name="location_name"
            value={profile.location_name || ''}
            onChange={handleChange}
          />
        </div>
        <button type="submit">Update Profile</button>
      </form>
    </div>
  );
};

export default Profile;