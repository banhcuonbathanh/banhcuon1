import { DishStatusValues } from '@/constants/type'
import z from 'zod'




export const CreateDishBody = z.object({
  name: z.string().min(1).max(256),
  price: z.coerce.number().positive(),
  description: z.string().max(10000),
  image: z.string().url(),
  status: z.enum(DishStatusValues).optional()
})

export type CreateDishBodyType = z.TypeOf<typeof CreateDishBody>

const DishSchema = z.object({
  id: z.number(),
  name: z.string(),
  price: z.coerce.number(),
  description: z.string(),
  image: z.string(),
  status: z.enum(DishStatusValues),
  createdAt: z.date(),
  updatedAt: z.date(),
  setId: z.number().optional() // New field to associate a dish with a set
});

export const DishRes = z.object({
  data: DishSchema,
  message: z.string()
})

export type DishResType = z.TypeOf<typeof DishRes>

export const DishListRes = z.array(DishSchema)

export type DishListResType = z.TypeOf<typeof DishListRes>

export const UpdateDishBody = CreateDishBody
export type UpdateDishBodyType = CreateDishBodyType
export const DishParams = z.object({
  id: z.coerce.number()
})
export type DishParamsType = z.TypeOf<typeof DishParams>
export type Dish = z.TypeOf<typeof DishSchema>;


/// set schema

const SetSchema = z.object({
  id: z.number(),
  name: z.string(),
  description: z.string().optional(),
  dishes: z.array(DishSchema)
});

export const SetListRes = z.array(SetSchema);
export type SetListResType = z.TypeOf<typeof SetListRes>;

export type Set = z.TypeOf<typeof SetSchema>;

// set favourite dish

const FavoriteSetSchema = z.object({
  id: z.number(),
  userId: z.number(), // Assuming you have user authentication
  name: z.string(),
  dishes: z.array(z.number()), // Array of dish IDs
  createdAt: z.date(),
  updatedAt: z.date()
});
export const FavoriteSetListRes = z.array(FavoriteSetSchema);
export type FavoriteSetListResType = z.TypeOf<typeof FavoriteSetListRes>;
export type FavoriteSet = z.TypeOf<typeof FavoriteSetSchema>;