import { z } from 'zod'

const configSchema = z.object({
  NEXT_PUBLIC_API_ENDPOINT: z.string(),
  NEXT_PUBLIC_URL: z.string(),
  NEXT_PUBLIC_API_Create_User:z.string(),
  NEXT_PUBLIC_API_Get_Account_Email:z.string(),

  NEXT_PUBLIC_API_Login:z.string()
})
// /users/email
const configProject = configSchema.safeParse({
  NEXT_PUBLIC_API_ENDPOINT: process.env.NEXT_PUBLIC_API_ENDPOINT,
  NEXT_PUBLIC_URL: process.env.NEXT_PUBLIC_URL,
  NEXT_PUBLIC_API_Create_User: process.env.NEXT_PUBLIC_API_Create_User,

  NEXT_PUBLIC_API_Get_Account_Email: process.env.NEXT_PUBLIC_API_Get_Account_Email,
  NEXT_PUBLIC_API_Login: process.env.NEXT_PUBLIC_API_Login
})

if (!configProject.success) {
  console.error(configProject.error.errors)
  throw new Error('Các khai báo biến môi trường không hợp lệ')
}

const envConfig = configProject.data

export default envConfig
