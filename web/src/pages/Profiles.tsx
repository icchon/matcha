import React, { useState, useEffect } from 'react';
import { Link } from 'react-router-dom';
import { getAllProfiles } from '../services/commonService';
import { UserProfile } from '../services/profileService';

const Profiles: React.FC = () => {
  const [profiles, setProfiles] = useState<UserProfile[]>([]);
  const [loading, setLoading] = useState(true);

  useEffect(() => {
    const fetchProfiles = async () => {
      try {
        const response = await getAllProfiles();
        setProfiles(response.data);
      } catch (error) {
        console.error('Error fetching profiles', error);
      } finally {
        setLoading(false);
      }
    };

    fetchProfiles();
  }, []);

  if (loading) {
    return <div>Loading...</div>;
  }

  return (
    <div>
      <h1>User Profiles</h1>
      <ul>
        {profiles.map((profile) => (
          <li key={profile.user_id}>
            <Link to={`/users/${profile.user_id}/profile`}>
              {profile.username} ({profile.first_name} {profile.last_name})
            </Link>
          </li>
        ))}
      </ul>
    </div>
  );
};

export default Profiles;
