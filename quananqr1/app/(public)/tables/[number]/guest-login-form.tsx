import { useEffect } from 'react'
import { useForm } from 'react-hook-form'
import { zodResolver } from '@hookform/resolvers/zod'
import { useSearchParams, useParams, useRouter } from 'next/navigation'
import { Card, CardHeader, CardTitle, CardContent } from '@/components/ui/card'
import { Form, FormField, FormItem, FormMessage } from '@/components/ui/form'
import { Input } from '@/components/ui/input'
import { Label } from '@/components/ui/label'
import { Button } from '@/components/ui/button'
import { Alert, AlertDescription } from '@/components/ui/alert'
import { Loader2 } from 'lucide-react'
import { handleErrorApi } from '@/lib/utils'
import { GuestLoginBodyType, GuestLoginBody } from '@/schemaValidations/guest.schema'
import { useAuthStore } from '@/zusstand/auth/controller/auth-controller'


export default function GuestLoginForm() {
  const { guestLogin, loading, error, clearError } = useAuthStore()
  const searchParams = useSearchParams()
  const params = useParams()
  const tableNumber = Number(params.number)
  const token = searchParams.get('token')
  const router = useRouter()

  const form = useForm<GuestLoginBodyType>({
    resolver: zodResolver(GuestLoginBody),
    defaultValues: {
      name: '',
      token: token ?? '',
      tableNumber
    }
  })

  useEffect(() => {
    if (!token) {
      router.push('/')
    }
  }, [token, router])

  useEffect(() => {
    // Clear any previous errors when the component mounts
    clearError()
  }, [clearError])

  async function onSubmit(values: GuestLoginBodyType) {
    try {
      await guestLogin(values)
      if (!error) {
        router.push('/guest/menu')
      }
    } catch (error) {
      handleErrorApi({
        error,
        setError: form.setError
      })
    }
  }

  return (
    <Card className='mx-auto max-w-sm'>
      <CardHeader>
        <CardTitle className='text-2xl'>Đăng nhập gọi món</CardTitle>
      </CardHeader>
      <CardContent>
        {error && (
          <Alert variant="destructive" className="mb-4">
            <AlertDescription>{error}</AlertDescription>
          </Alert>
        )}
        <Form {...form}>
          <form
            className='space-y-2 max-w-[600px] flex-shrink-0 w-full'
            noValidate
            onSubmit={form.handleSubmit(onSubmit)}
          >
            <div className='grid gap-4'>
              <FormField
                control={form.control}
                name='name'
                render={({ field }) => (
                  <FormItem>
                    <div className='grid gap-2'>
                      <Label htmlFor='name'>Tên khách hàng</Label>
                      <Input id='name' type='text' required {...field} disabled={loading} />
                      <FormMessage />
                    </div>
                  </FormItem>
                )}
              />

              <Button type='submit' className='w-full' disabled={loading}>
                {loading ? (
                  <>
                    <Loader2 className="mr-2 h-4 w-4 animate-spin" />
                    Đang xử lý...
                  </>
                ) : (
                  'Đăng nhập'
                )}
              </Button>
            </div>
          </form>
        </Form>
      </CardContent>
    </Card>
  )
}