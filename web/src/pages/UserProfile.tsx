import React, { useState, useEffect } from 'react';
import { useParams } from 'react-router-dom';
import { getUserProfile, getUserPictures, UserProfile as UserProfileType, Picture } from '../services/profileService';
import { likeUser, unlikeUser, blockUser } from '../services/userService';

const UserProfile: React.FC = () => {
  const { userId } = useParams<{ userId: string }>();
  const [profile, setProfile] = useState<UserProfileType | null>(null);
  const [pictures, setPictures] = useState<Picture[]>([]);
  const [loading, setLoading] = useState(true);

  useEffect(() => {
    const fetchProfile = async () => {
      if (!userId) return;
      try {
        const response = await getUserProfile(userId);
        setProfile(response.data);
      } catch (error) {
        console.error('Error fetching profile', error);
      }
    };

    const fetchPictures = async () => {
        if (!userId) return;
        try {
            const response = await getUserPictures(userId);
            setPictures(response.data);
        } catch (error) {
            console.error('Error fetching pictures', error);
        }
    }

    const loadData = async () => {
        setLoading(true);
        await Promise.all([fetchProfile(), fetchPictures()]);
        setLoading(false);
    }

    loadData();
  }, [userId]);

  const handleLike = async () => {
    if (!userId) return;
    try {
      await likeUser(userId);
      alert('User liked');
    } catch (error) {
      alert('Failed to like user');
    }
  };

  const handleUnlike = async () => {
    if (!userId) return;
    try {
      await unlikeUser(userId);
      alert('User unliked');
    } catch (error) {
      alert('Failed to unlike user');
    }
  };

  const handleBlock = async () => {
    if (!userId) return;
    try {
      await blockUser(userId);
      alert('User blocked');
    } catch (error) {
      alert('Failed to block user');
    }
  };


  if (loading) {
    return <div>Loading...</div>;
  }

  if (!profile) {
    return <div>Failed to load profile.</div>;
  }

  return (
    <div>
      <h1>{profile.username}'s Profile</h1>
      <div style={{ display: 'flex', flexWrap: 'wrap', gap: '10px', marginBottom: '20px' }}>
          {pictures.map(pic => (
              <div key={pic.id} style={{ border: pic.is_profile_pic ? '2px solid blue' : 'none', padding: '5px' }}>
                  <img src={pic.url} alt="User" width="150" />
              </div>
          ))}
      </div>
      <p>
        <strong>First Name:</strong> {profile.first_name}
      </p>
      <p>
        <strong>Last Name:</strong> {profile.last_name}
      </p>
      <p>
        <strong>Gender:</strong> {profile.gender}
      </p>
      <p>
        <strong>Sexual Preference:</strong> {profile.sexual_preference}
      </p>
      <p>
        <strong>Birthday:</strong> {profile.birthday}
      </p>
      <p>
        <strong>Occupation:</strong> {profile.occupation}
      </p>
      <p>
        <strong>Biography:</strong> {profile.biography}
      </p>
      <p>
        <strong>Location:</strong> {profile.location_name}
      </p>
      <button onClick={handleLike}>Like</button>
      <button onClick={handleUnlike}>Unlike</button>
      <button onClick={handleBlock}>Block</button>
    </div>
  );
};

export default UserProfile;