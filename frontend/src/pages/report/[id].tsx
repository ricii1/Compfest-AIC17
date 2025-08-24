import axios from 'axios';
import { useRouter } from 'next/router';
import * as React from 'react';
import { useEffect, useState } from 'react';

import { getFromSessionStorage } from '@/lib/helper';

import Layout from '@/components/layout/Layout';
import Seo from '@/components/Seo';
import { Badge } from '@/components/ui/badge';
import {
  Select,
  SelectContent,
  SelectItem,
  SelectTrigger,
  SelectValue,
} from '@/components/ui/select';

import { backendUrl } from '@/constant/env';

type StatusVariant = 'verified' | 'rejected' | 'handled' | 'unverified';

interface ReportData {
  id: string;
  text: string;
  image: string;
  location: string;
  status: StatusVariant;
  upvotes: number;
  share_count: number;
  tag_id: string;
  user_id: string;
  username: string;
  pred_confidence: number;
  created_at: string;
  tag: Tag;
}

interface Tag {
  class: string;
  id: string;
  location: string;
}

interface ApiResponse {
  status: boolean;
  message: string;
  data: ReportData;
}

export default function ReportDetailPage() {
  const router = useRouter();
  const { id } = router.query;
  const [status, setStatus] = useState<StatusVariant>('unverified');
  const [report, setReport] = useState<ReportData | null>(null);
  const [loading, setLoading] = useState(true);
  const [error, setError] = useState<string | null>(null);

  // Helper function to handle auth errors
  const handleAuthError = (error: any) => {
    if (axios.isAxiosError(error)) {
      // Check for 401 (Unauthorized) or 403 (Forbidden)
      if (error.response?.status === 401 || error.response?.status === 403) {
        console.log('Authentication error, redirecting to login...');
        // Clear token from storage
        if (typeof window !== 'undefined') {
          sessionStorage.removeItem('access_token');
          localStorage.removeItem('access_token');
        }
        // Redirect to login
        router.push('/login');
        return true; // Indicates auth error was handled
      }
      // Check for 400 Bad Request that might indicate token issues
      if (error.response?.status === 400) {
        const errorMessage = error.response?.data?.message || '';
        if (errorMessage.toLowerCase().includes('token') ||
          errorMessage.toLowerCase().includes('unauthorized') ||
          errorMessage.toLowerCase().includes('invalid')) {
          console.log('Token-related bad request, redirecting to login...');
          // Clear token from storage
          if (typeof window !== 'undefined') {
            sessionStorage.removeItem('access_token');
            localStorage.removeItem('access_token');
          }
          // Redirect to login
          router.push('/login');
          return true; // Indicates auth error was handled
        }
      }
    }
    return false; // Auth error was not handled
  };

  // Fetch report data from API
  useEffect(() => {
    const fetchReport = async () => {
      if (!id) return;

      try {
        setLoading(true);
        setError(null);
        const token = getFromSessionStorage('access_token');
        const headers = token ? { Authorization: `Bearer ${token}` } : {};

        const response = await axios.get<ApiResponse>(
          `${backendUrl}/api/reports/${id}`,
          { headers }
        );

        if (response.status === 200 && response.data.status) {
          setReport(response.data.data);
          setStatus(response.data.data.status);
        } else {
          throw new Error(response.data.message || 'Failed to fetch report');
        }
      } catch (error) {
        console.error('Error fetching report:', error);
        // Check if it's an auth error and handle redirect
        if (!handleAuthError(error)) {
          // If not auth error, show general error message
          setError('Failed to load report. Please try again later.');
        }
      } finally {
        setLoading(false);
      }
    };

    fetchReport();
  }, [id]);

  const [submitting, setSubmitting] = useState(false);

  const handleStatusChange = (newStatus: string) => {
    setStatus(newStatus as StatusVariant);
  };

  const handleSubmitStatus = async () => {
    if (!id || !report) return;

    try {
      setSubmitting(true);
      const token = getFromSessionStorage('access_token');
      const headers = token ? { Authorization: `Bearer ${token}` } : {};

      // Send API request to update status
      const response = await axios.post(
        `${backendUrl}/api/reports/${id}/status`,
        { status },
        { headers }
      );

      if (response.status === 200) {
        // Update local report state
        setReport(prev => prev ? { ...prev, status } : prev);
      } else {
        throw new Error('Failed to update status');
      }
    } catch (error) {
      console.error('Error updating status:', error);
      // Check if it's an auth error and handle redirect
      if (!handleAuthError(error)) {
        // If not auth error, show general error message
        alert(`Failed to update status: ${error instanceof Error ? error.message : 'Unknown error'}`);
      }
    } finally {
      setSubmitting(false);
    }
  };

  if (loading) {
    return (
      <Layout>
        <Seo templateTitle="Loading Report..." />
        <main className='min-h-screen bg-gray-50 p-6'>
          <div className='mx-auto max-w-4xl'>
            <div className="animate-pulse">
              <div className="h-8 bg-gray-300 rounded w-1/3 mb-6"></div>
              <div className="bg-white p-6 rounded-lg shadow">
                <div className="h-64 bg-gray-300 rounded mb-6"></div>
                <div className="h-4 bg-gray-300 rounded w-full mb-4"></div>
                <div className="h-4 bg-gray-300 rounded w-3/4"></div>
              </div>
            </div>
          </div>
        </main>
      </Layout>
    );
  }

  if (error || !report) {
    return (
      <Layout>
        <Seo templateTitle="Error - Report Detail" />
        <main className='min-h-screen bg-gray-50 p-6'>
          <div className='mx-auto max-w-4xl'>
            <div className='bg-red-50 border border-red-200 rounded-lg p-6 text-center'>
              <div className="text-red-600 text-xl mb-2">⚠️ Error</div>
              <p className="text-red-800">{error || 'Report not found'}</p>
              <button
                onClick={() => router.back()}
                className="mt-4 px-4 py-2 bg-blue-600 text-white rounded-md hover:bg-blue-700 transition"
              >
                Go Back
              </button>
            </div>
          </div>
        </main>
      </Layout>
    );
  }

  return (
    <Layout>
      <Seo templateTitle={`Report Detail - ${report.id}`} />
      <main className='min-h-screen bg-gray-50 p-6'>
        <div className='mx-auto max-w-4xl'>
          <div className='flex items-center justify-between mb-6'>
            <h1 className='text-3xl font-bold text-gray-900'>Report Detail</h1>
            <button
              onClick={() => router.back()}
              className='px-4 py-2 text-sm font-medium text-gray-700 bg-white border border-gray-300 rounded-md hover:bg-gray-50 transition'
            >
              ← Go Back
            </button>
          </div>

          <div className='bg-white shadow rounded-lg p-6'>
            {/* Image */}
            {report.image && (
              <div className='mb-6'>
                <h2 className='text-lg font-semibold text-gray-800 mb-2'>Image</h2>
                <img
                  src={report.image.startsWith('http') ? report.image : `${backendUrl}/${report.image}`}
                  alt='Report Image'
                  className='w-full h-auto rounded-lg max-h-96 object-cover'
                />
              </div>
            )}

            {/* Text */}
            <div className='mb-4'>
              <h2 className='text-lg font-semibold text-gray-800'>Description</h2>
              <p className='text-gray-600 mt-2'>{report.text}</p>
            </div>

            {/* User */}
            <div className='mb-4'>
              <h2 className='text-lg font-semibold text-gray-800'>Reported By</h2>
              <p className='text-gray-600 mt-2'>
                {report.username} (ID: {report.user_id})
              </p>
            </div>

            {/* Location */}
            {report.location && (
              <div className='mb-4'>
                <h2 className='text-lg font-semibold text-gray-800'>Location</h2>
                <p className='text-gray-600 mt-2'>{report.location}</p>
              </div>
            )}

            {/* Created At */}
            <div className='mb-4'>
              <h2 className='text-lg font-semibold text-gray-800'>Created At</h2>
              <p className='text-gray-600 mt-2'>
                {report.created_at ?
                  new Date(report.created_at).toLocaleString() :
                  'No date available'
                }
              </p>
            </div>

            {/* Confidence */}
            <div className='mb-4'>
              <h2 className='text-lg font-semibold text-gray-800'>AI Confidence</h2>
              <div className='flex items-center mt-2'>
                <span className='text-sm font-medium'>{report.pred_confidence}%</span>
                <div className='ml-2 w-24 bg-gray-200 rounded-full h-2'>
                  <div
                    className='bg-blue-600 h-2 rounded-full'
                    style={{ width: `${Math.max(report.pred_confidence, 5)}%` }}
                  ></div>
                </div>
              </div>
            </div>

            {/* Tag */}
            {report.tag_id && (
              <div className='mb-4'>
                <h2 className='text-lg font-semibold text-gray-800'>Tag</h2>
                {report.tag.class.split(',').map((tag, index) => (
                  <Badge key={index} variant="outline">
                    {tag.trim()}
                  </Badge>
                ))}
              </div>
            )}

            {/* Status Dropdown */}
            <div className='mb-4'>
              <h2 className='text-lg font-semibold text-gray-800'>Status</h2>
              <div className='mt-2'>
                <Select value={status} onValueChange={handleStatusChange}>
                  <SelectTrigger>
                    <SelectValue placeholder='Select a status' />
                  </SelectTrigger>
                  <SelectContent>
                    <SelectItem value='unverified'>Unverified</SelectItem>
                    <SelectItem value='verified'>Verified</SelectItem>
                    <SelectItem value='handled'>Handled</SelectItem>
                    <SelectItem value='rejected'>Rejected</SelectItem>
                  </SelectContent>
                </Select>
              </div>
            </div>

            {/* Current Status Badge */}
            <div className='mt-6'>
              <h2 className='text-lg font-semibold text-gray-800'>Current Status</h2>
              <Badge variant={report.status}>{report.status}</Badge>
            </div>

            {/* Submit Button */}
            <div className='mt-6'>
              <button
                onClick={handleSubmitStatus}
                disabled={submitting || status === report.status}
                className={`px-4 py-2 rounded-md transition ${submitting || status === report.status
                  ? 'bg-gray-400 cursor-not-allowed'
                  : 'bg-blue-600 hover:bg-blue-700'
                  } text-white`}
              >
                {submitting ? 'Updating...' : status === report.status ? 'No Changes' : 'Update Status'}
              </button>
            </div>
          </div>
        </div>
      </main>
    </Layout>
  );
}