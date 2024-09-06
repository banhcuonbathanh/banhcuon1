// src/pages/index.tsx (SSG example)
import { GetStaticProps } from 'next';
import { User } from './domain/entities/user';
import { fetchUser } from '@/app/lib/api';


interface HomeProps {
  user: User;
}

export default function Home({ user }: HomeProps) {
  return (
    <div>
      <h1>Welcome, {user.name}!</h1>
      <p>Email: {user.email}</p>
    </div>
  );
}

export const getStaticProps: GetStaticProps<HomeProps> = async () => {
  const user = await fetchUser(1); // Fetch a default user
  return { props: { user }, revalidate: 60 };
};