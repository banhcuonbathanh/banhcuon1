import { fetchUser } from '@/app/lib/api';
import { NextApiRequest, NextApiResponse } from 'next';
``

export default async function handler(req: NextApiRequest, res: NextApiResponse) {
  const { id } = req.query;
  try {
    const user = await fetchUser(Number(id));
    res.status(200).json(user);
  } catch (error) {
    res.status(500).json({ error: 'Failed to fetch user' });
  }
}