// components/UserList.tsx
"use client";
import React, { useEffect } from "react";
import { useUserStore } from "./controller/controller_user";

const UserList: React.FC = () => {
  const { users, loading, error, fetchUsers } = useUserStore();

  useEffect(() => {
    fetchUsers();
  }, [fetchUsers]);

  if (loading) return <div>Loading...</div>;
  if (error) return <div>Error: {error}</div>;

  return (
    <ul>
      {users.map((user) => (
        <li key={user.id}>
          {user.name} ({user.email})
        </li>
      ))}
    </ul>
  );
};

export default UserList;
