import { fetchDataDelay } from '@/zusstand/auth/server/server-auth-controler';

import React from 'react'

const ManageHomePage = async () => {

  const account = await fetchDataDelay();
  throw new Error('Failed to fetch account data');
  return (
    <div>page</div>
  )
}

export default ManageHomePage