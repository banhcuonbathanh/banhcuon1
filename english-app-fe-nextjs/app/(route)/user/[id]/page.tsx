// src/pages/users/[id].tsx (Server-side fetching example)
import { GetServerSideProps } from 'next';
import { User } from '../domain/entities/user';
import { fetchUser } from '@/app/lib/api';


interface UserPageProps {
  user: User;
}

export default function UserPage({ user }: UserPageProps) {
  return (
    <div>
      <h1>{user.name}</h1>
      <p>Email: {user.email}</p>
    </div>
  );
}

export const getServerSideProps: GetServerSideProps<UserPageProps> = async (context) => {
  const id = Number(context.params?.id);
  const user = await fetchUser(id);
  return { props: { user } };
};