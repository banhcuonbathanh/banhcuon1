// src/components/UserProfile.tsx (Client-side fetching example)
import { useState, useEffect } from 'react';
import { User } from '../domain/entities/user';

interface UserProfileProps {
  userId: number;
}

export default function UserProfile({ userId }: UserProfileProps) {
  const [user, setUser] = useState<User | null>(null);

  useEffect(() => {
    async function fetchUser() {
      const response = await fetch(`/api/users/${userId}`);
      const data = await response.json();
      setUser(data);
    }
    fetchUser();
  }, [userId]);

  if (!user) return <div>Loading...</div>;

  return (
    <div>
      <h2>{user.name}</h2>
      <p>Email: {user.email}</p>
    </div>
  );
}
