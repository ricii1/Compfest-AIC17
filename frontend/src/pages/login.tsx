import Layout from '@/components/layout/Layout';
import Seo from '@/components/Seo';
import { Button } from '@/components/ui/button';
import {
  Card,
  CardContent,
  CardDescription,
  CardFooter,
  CardHeader,
  CardTitle,
} from '@/components/ui/card';
import { Input } from '@/components/ui/input';
import { Label } from '@/components/ui/label';
import { useForm } from 'react-hook-form';
import { backendUrl } from '@/constant/env';
import axios from "axios";
import Swal from 'sweetalert2';
import { setToSessionStorage } from '@/lib/helper';

type LoginFormData = {
  email: string;
  password: string;
};

export default function LoginPage() {
  const {
    register,
    handleSubmit,
    formState: { errors, isSubmitting },
  } = useForm<LoginFormData>();

  const onSubmit = async (data: LoginFormData) => {
    try {
      // Handle login logic here
      handleLogin(data);
    } catch (error) {
      console.error('Login failed:', error);
    }
  };

  const handleLogin = async (data: LoginFormData) => {
    try {
      const response = await axios.post(`${backendUrl}/api/user/login`, data);

      if (response.status !== 200) {
        throw new Error('Login failed');
      }
      const result = response.data.data;
      if (result.role === 'admin') {
        await Swal.fire({
          icon: 'success',
          title: 'Login Successful!',
          text: 'Welcome to RAPID Admin Dashboard',
          timer: 2000,
          showConfirmButton: false
        });
        setToSessionStorage("access_token", result.access_token)
        window.location.href = '/';
      } else if (result.role === 'user') {
        await Swal.fire({
          icon: 'error',
          title: 'Access Denied',
          text: 'You do not have admin privileges.',
          confirmButtonColor: '#d33'
        });
      } else {
        await Swal.fire({
          icon: 'error',
          title: 'Login Failed',
          text: 'Unknown role.',
          confirmButtonColor: '#d33'
        });
      }
    } catch (error) {
      await Swal.fire({
        icon: 'error',
        title: 'Login Failed',
        text: 'Please check your credentials and try again.',
        confirmButtonColor: '#d33'
      });
    }
  };

  return (
    <Layout>
      <Seo templateTitle='Login' />
      <main className='flex min-h-screen items-center justify-center bg-gradient-to-br from-blue-50 via-white to-blue-100 p-4'>
        <div className='w-full max-w-md'>
          {/* RAPID Logo/Branding */}
          <div className='mb-8 text-center'>
            <h1 className='text-3xl font-bold text-gray-900'>RAPID</h1>
            <p className='text-sm text-gray-600'>Real-Time AI Problem and Incident Dispatch</p>
          </div>

          {/* Login Card */}
          <Card className='shadow-xl border-0 bg-white/80 backdrop-blur-sm'>
            <CardHeader className='pb-4 pt-8 px-8'>
              <CardTitle className='text-2xl font-semibold text-center text-gray-900'>
                Welcome Back
              </CardTitle>
              <CardDescription className='text-center text-gray-600'>
                Sign in to your RAPID account to continue
              </CardDescription>
            </CardHeader>

            <CardContent className='px-8'>
              <form onSubmit={handleSubmit(onSubmit)} className='space-y-6'>
                <div className='space-y-2'>
                  <Label
                    htmlFor='email'
                    className='text-sm font-medium text-gray-700'
                  >
                    Email Address
                  </Label>
                  <Input
                    id='email'
                    type='email'
                    placeholder='Enter your email'
                    {...register('email', {
                      required: 'Email is required',
                      pattern: {
                        value: /^[A-Z0-9._%+-]+@[A-Z0-9.-]+\.[A-Z]{2,}$/i,
                        message: 'Invalid email address',
                      },
                    })}
                    className={`h-11 border-gray-300 focus:border-blue-500 focus:ring-blue-500 ${errors.email ? 'border-red-500 focus:border-red-500 focus:ring-red-500' : ''
                      }`}
                  />
                  {errors.email && (
                    <p className='text-sm text-red-600'>{errors.email.message}</p>
                  )}
                </div>

                <div className='space-y-2'>
                  <div className='flex items-center justify-between'>
                    <Label
                      htmlFor='password'
                      className='text-sm font-medium text-gray-700'
                    >
                      Password
                    </Label>
                  </div>
                  <Input
                    id='password'
                    type='password'
                    placeholder='Enter your password'
                    {...register('password', {
                      required: 'Password is required',
                      minLength: {
                        value: 6,
                        message: 'Password must be at least 6 characters',
                      },
                    })}
                    className={`h-11 border-gray-300 focus:border-blue-500 focus:ring-blue-500 ${errors.password ? 'border-red-500 focus:border-red-500 focus:ring-red-500' : ''
                      }`}
                  />
                  {errors.password && (
                    <p className='text-sm text-red-600'>{errors.password.message}</p>
                  )}
                </div>

                <Button
                  type='submit'
                  disabled={isSubmitting}
                  className='w-full h-11 bg-gradient-to-r from-blue-600 to-blue-700 hover:from-blue-700 hover:to-blue-800 text-white font-medium transition-all duration-200 shadow-lg hover:shadow-xl disabled:opacity-50 disabled:cursor-not-allowed'
                >
                  {isSubmitting ? 'Signing In...' : 'Sign In'}
                </Button>
              </form>
            </CardContent>
          </Card>

          {/* Footer */}
          <div className='mt-8 text-center text-xs text-gray-500'>
            <p>&copy; 2025 RAPID. All rights reserved.</p>
          </div>
        </div>
      </main>
    </Layout>
  );
}
