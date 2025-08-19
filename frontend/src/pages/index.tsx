import * as React from 'react';
import { useState, useEffect } from 'react';

import Layout from '@/components/layout/Layout';
import Seo from '@/components/Seo';
import { Button } from '@/components/ui/button';
import { Badge } from '@/components/ui/badge';
import {
  Table,
  TableHeader,
  TableBody,
  TableHead,
  TableRow,
  TableCell,
} from '@/components/ui/table';
import { Card, CardContent, CardHeader, CardTitle } from '@/components/ui/card';

// Type definitions based on the Report struct
type ReportStatus = 'verified' | 'unverified' | 'rejected';

interface User {
  id: string;
  name: string;
  email: string;
}

interface Tag {
  id: string;
  name: string;
  color: string;
}

interface Report {
  id: string;
  text: string;
  image: string;
  status: ReportStatus;
  pred_confidence: number | null;
  upvotes: number;
  share_count: number;
  location: string;
  user_id: string;
  user: User;
  tag_id: string | null;
  tag: Tag | null;
  created_at: string;
  updated_at: string;
}

export default function AdminDashboard() {
  const [reports, setReports] = useState<Report[]>([]);
  const [loading, setLoading] = useState(true);

  // Mock data for demonstration
  useEffect(() => {
    const mockReports: Report[] = [
      {
        id: '1',
        text: 'Jalan rusak di Jl. Sudirman',
        image: '/images/road-damage.jpg',
        status: 'unverified',
        pred_confidence: 85,
        upvotes: 12,
        share_count: 3,
        location: 'Jl. Sudirman, Jakarta',
        user_id: 'user1',
        user: { id: 'user1', name: 'John Doe', email: 'john@example.com' },
        tag_id: 'tag1',
        tag: { id: 'tag1', name: 'Infrastruktur', color: '#3B82F6' },
        created_at: '2025-01-15T10:30:00Z',
        updated_at: '2025-01-15T10:30:00Z',
      },
      {
        id: '2',
        text: 'Lampu jalan mati di perempatan',
        image: '/images/streetlight.jpg',
        status: 'verified',
        pred_confidence: 92,
        upvotes: 8,
        share_count: 1,
        location: 'Perempatan Gatot Subroto',
        user_id: 'user2',
        user: { id: 'user2', name: 'Jane Smith', email: 'jane@example.com' },
        tag_id: 'tag2',
        tag: { id: 'tag2', name: 'Penerangan', color: '#F59E0B' },
        created_at: '2025-01-14T15:45:00Z',
        updated_at: '2025-01-15T09:20:00Z',
      },
      {
        id: '3',
        text: 'Sampah berserakan di taman kota',
        image: '/images/garbage.jpg',
        status: 'rejected',
        pred_confidence: 67,
        upvotes: 5,
        share_count: 0,
        location: 'Taman Menteng, Jakarta',
        user_id: 'user3',
        user: { id: 'user3', name: 'Bob Wilson', email: 'bob@example.com' },
        tag_id: 'tag3',
        tag: { id: 'tag3', name: 'Kebersihan', color: '#10B981' },
        created_at: '2025-01-13T08:15:00Z',
        updated_at: '2025-01-14T11:30:00Z',
      },
    ];

    // Simulate API call
    setTimeout(() => {
      setReports(mockReports);
      setLoading(false);
    }, 1000);
  }, []);

  const getStatusBadge = (status: ReportStatus) => {
    switch (status) {
      case 'verified':
        return <Badge variant="verified">Verified</Badge>;
      case 'unverified':
        return <Badge variant="unverified">Unverified</Badge>;
      case 'rejected':
        return <Badge variant="rejected">Rejected</Badge>;
      default:
        return <Badge variant="secondary">{status}</Badge>;
    }
  };

  const handleStatusChange = (reportId: string, newStatus: ReportStatus) => {
    setReports(prev =>
      prev.map(report =>
        report.id === reportId
          ? { ...report, status: newStatus, updated_at: new Date().toISOString() }
          : report
      )
    );
  };

  if (loading) {
    return (
      <Layout>
        <Seo templateTitle="Admin Dashboard - RAPID" />
        <main className="min-h-screen bg-gray-50 p-6">
          <div className="mx-auto max-w-7xl">
            <div className="animate-pulse">
              <div className="h-8 bg-gray-300 rounded w-1/3 mb-6"></div>
              <div className="bg-white p-6 rounded-lg shadow">
                <div className="h-4 bg-gray-300 rounded w-full mb-4"></div>
                <div className="h-4 bg-gray-300 rounded w-3/4 mb-4"></div>
                <div className="h-4 bg-gray-300 rounded w-1/2"></div>
              </div>
            </div>
          </div>
        </main>
      </Layout>
    );
  }
  return (
    <Layout>
      <Seo templateTitle="Admin Dashboard - RAPID" />
      <main className="min-h-screen bg-gray-50 p-6">
        <div className="mx-auto max-w-7xl">
          {/* Header */}
          <div className="mb-8">
            <div className="flex items-center justify-between">
              <div>
                <h1 className="text-3xl font-bold text-gray-900">RAPID Admin Dashboard</h1>
                <p className="text-gray-600 mt-2">Manage and monitor user reports</p>
              </div>
              <div className="flex items-center space-x-2">
                <span className="text-sm text-gray-500">Total Reports:</span>
                <Badge variant="info">{reports.length}</Badge>
              </div>
            </div>
          </div>

          {/* Stats Cards */}
          <div className="grid grid-cols-1 md:grid-cols-4 gap-6 mb-8">
            <Card>
              <CardHeader className="pb-2">
                <CardTitle className="text-sm font-medium text-gray-600">Total Reports</CardTitle>
              </CardHeader>
              <CardContent>
                <div className="text-2xl font-bold text-gray-900">{reports.length}</div>
              </CardContent>
            </Card>
            <Card>
              <CardHeader className="pb-2">
                <CardTitle className="text-sm font-medium text-gray-600">Verified</CardTitle>
              </CardHeader>
              <CardContent>
                <div className="text-2xl font-bold text-green-600">
                  {reports.filter(r => r.status === 'verified').length}
                </div>
              </CardContent>
            </Card>
            <Card>
              <CardHeader className="pb-2">
                <CardTitle className="text-sm font-medium text-gray-600">Unverified</CardTitle>
              </CardHeader>
              <CardContent>
                <div className="text-2xl font-bold text-yellow-600">
                  {reports.filter(r => r.status === 'unverified').length}
                </div>
              </CardContent>
            </Card>
            <Card>
              <CardHeader className="pb-2">
                <CardTitle className="text-sm font-medium text-gray-600">Rejected</CardTitle>
              </CardHeader>
              <CardContent>
                <div className="text-2xl font-bold text-red-600">
                  {reports.filter(r => r.status === 'rejected').length}
                </div>
              </CardContent>
            </Card>
          </div>

          {/* Reports Table */}
          <Card>
            <CardHeader>
              <CardTitle>Reports Management</CardTitle>
            </CardHeader>
            <CardContent>
              <div className="overflow-x-auto">
                <Table>
                  <TableHeader>
                    <TableRow>
                      <TableHead className="w-20">ID</TableHead>
                      <TableHead className="min-w-[200px]">Report Text</TableHead>
                      <TableHead>Status</TableHead>
                      <TableHead>Confidence</TableHead>
                      <TableHead>Upvotes</TableHead>
                      <TableHead>Shares</TableHead>
                      <TableHead className="min-w-[150px]">Location</TableHead>
                      <TableHead>User</TableHead>
                      <TableHead>Tag</TableHead>
                      <TableHead>Created</TableHead>
                      <TableHead className="w-[200px]">Actions</TableHead>
                    </TableRow>
                  </TableHeader>
                  <TableBody>
                    {reports.map((report) => (
                      <TableRow key={report.id}>
                        <TableCell className="font-mono text-xs">
                          {report.id.slice(0, 8)}...
                        </TableCell>
                        <TableCell>
                          <div className="max-w-xs">
                            <p className="truncate font-medium">{report.text}</p>
                            {report.image && (
                              <p className="text-xs text-gray-500 mt-1">📷 Has image</p>
                            )}
                          </div>
                        </TableCell>
                        <TableCell>
                          {getStatusBadge(report.status)}
                        </TableCell>
                        <TableCell>
                          {report.pred_confidence ? (
                            <div className="flex items-center">
                              <span className="text-sm font-medium">{report.pred_confidence}%</span>
                              <div className="ml-2 w-12 bg-gray-200 rounded-full h-2">
                                <div
                                  className="bg-blue-600 h-2 rounded-full"
                                  style={{ width: `${report.pred_confidence}%` }}
                                ></div>
                              </div>
                            </div>
                          ) : (
                            <span className="text-gray-400">-</span>
                          )}
                        </TableCell>
                        <TableCell>
                          <div className="flex items-center">
                            <span className="text-sm">👍 {report.upvotes}</span>
                          </div>
                        </TableCell>
                        <TableCell>
                          <div className="flex items-center">
                            <span className="text-sm">🔗 {report.share_count}</span>
                          </div>
                        </TableCell>
                        <TableCell>
                          <span className="text-sm text-gray-600">{report.location}</span>
                        </TableCell>
                        <TableCell>
                          <div>
                            <p className="text-sm font-medium">{report.user.name}</p>
                            <p className="text-xs text-gray-500">{report.user.email}</p>
                          </div>
                        </TableCell>
                        <TableCell>
                          {report.tag ? (
                            <Badge
                              variant="outline"
                              style={{ borderColor: report.tag.color, color: report.tag.color }}
                            >
                              {report.tag.name}
                            </Badge>
                          ) : (
                            <span className="text-gray-400">-</span>
                          )}
                        </TableCell>
                        <TableCell>
                          <span className="text-xs text-gray-500">
                            {new Date(report.created_at).toLocaleDateString()}
                          </span>
                        </TableCell>
                        <TableCell>
                          <div className="flex space-x-1">
                            <Button
                              size="sm"
                              variant="outline"
                              className="text-green-600 border-green-600 hover:bg-green-50"
                              onClick={() => handleStatusChange(report.id, 'verified')}
                              disabled={report.status === 'verified'}
                            >
                              ✓
                            </Button>
                            <Button
                              size="sm"
                              variant="outline"
                              className="text-red-600 border-red-600 hover:bg-red-50"
                              onClick={() => handleStatusChange(report.id, 'rejected')}
                              disabled={report.status === 'rejected'}
                            >
                              ✕
                            </Button>
                            <Button
                              size="sm"
                              variant="outline"
                              className="text-yellow-600 border-yellow-600 hover:bg-yellow-50"
                              onClick={() => handleStatusChange(report.id, 'unverified')}
                              disabled={report.status === 'unverified'}
                            >
                              ?
                            </Button>
                          </div>
                        </TableCell>
                      </TableRow>
                    ))}
                  </TableBody>
                </Table>
              </div>
            </CardContent>
          </Card>
        </div>
      </main>
    </Layout>
  );
}
