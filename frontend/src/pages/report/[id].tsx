import * as React from 'react';
import { useState } from 'react';
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

type StatusVariant = 'verified' | 'rejected' | 'handled' | 'unverified';

export default function ReportDetailPage() {
  const [status, setStatus] = useState<StatusVariant>('unverified');

  // Mock data for demonstration
  const report = {
    id: '1',
    text: 'Jalan rusak di Jl. Sudirman',
    image: '/images/road-damage.jpg',
    status: 'unverified',
    pred_confidence: 85,
    location: 'Jl. Sudirman, Jakarta',
    user: { name: 'John Doe', email: 'john@example.com' },
    created_at: '2025-01-15T10:30:00Z',
  };

  const handleStatusChange = (newStatus: string) => {
    setStatus(newStatus as StatusVariant);
    // Simulate API call to update status
    console.log(`Status updated to: ${newStatus}`);
  };

  return (
    <Layout>
      <Seo templateTitle={`Report Detail - ${report.id}`} />
      <main className='min-h-screen bg-gray-50 p-6'>
        <div className='mx-auto max-w-4xl'>
          <h1 className='text-3xl font-bold text-gray-900 mb-6'>Report Detail</h1>

          <div className='bg-white shadow rounded-lg p-6'>
            {/* Image */}
            {report.image && (
              <div className='mb-6'>
                <img
                  src={report.image}
                  alt='Report Image'
                  className='w-full h-auto rounded-lg'
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
                {report.user.name} ({report.user.email})
              </p>
            </div>

            {/* Created At */}
            <div className='mb-4'>
              <h2 className='text-lg font-semibold text-gray-800'>Created At</h2>
              <p className='text-gray-600 mt-2'>
                {new Date(report.created_at).toLocaleString()}
              </p>
            </div>

            {/* Confidence */}
            <div className='mb-4'>
              <h2 className='text-lg font-semibold text-gray-800'>Confidence</h2>
              <div className='flex items-center mt-2'>
                <span className='text-sm font-medium'>{report.pred_confidence}%</span>
                <div className='ml-2 w-24 bg-gray-200 rounded-full h-2'>
                  <div
                    className='bg-blue-600 h-2 rounded-full'
                    style={{ width: `${report.pred_confidence}%` }}
                  ></div>
                </div>
              </div>
            </div>

            {/* Status Dropdown */}
            <div className='mb-4'>
              <h2 className='text-lg font-semibold text-gray-800'>Status</h2>
              <div className='mt-2'>
                <Select value={status} onValueChange={handleStatusChange}>
                  <SelectTrigger>
                    <SelectValue placeholder='Select a status' />
                  </SelectTrigger>
                  <SelectContent>
                    <SelectItem value='verified'>Verified</SelectItem>
                    <SelectItem value='rejected'>Rejected</SelectItem>
                    <SelectItem value='handled'>Handled</SelectItem>
                  </SelectContent>
                </Select>
              </div>
            </div>

            {/* Current Status Badge */}
            <div className='mt-6'>
              <h2 className='text-lg font-semibold text-gray-800'>Current Status</h2>
              <Badge variant={status}>{status}</Badge>
            </div>

            {/* Submit Button */}
            <div className='mt-6'>
              <button
                onClick={() => console.log(`Submitting status: ${status}`)}
                className='px-4 py-2 bg-blue-600 text-white rounded-md hover:bg-blue-700 transition'
              >
                Submit Status
              </button>
            </div>
          </div>
        </div>
      </main>
    </Layout>
  );
}