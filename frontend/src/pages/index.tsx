import axios from 'axios';
import { useRouter } from 'next/router';
import * as React from 'react';
import { useEffect, useState } from 'react';

import { getFromSessionStorage } from '@/lib/helper';

import Layout from '@/components/layout/Layout';
import Seo from '@/components/Seo';
import { Badge } from '@/components/ui/badge';
import { Button } from '@/components/ui/button';
import { Card, CardContent, CardHeader, CardTitle } from '@/components/ui/card';
import { Input } from '@/components/ui/input';
import {
  Table,
  TableBody,
  TableCell,
  TableHead,
  TableHeader,
  TableRow,
} from '@/components/ui/table';

import { backendUrl } from '@/constant/env';

type ReportStatus = 'unverified' | 'verified' | 'handled' | 'rejected';

interface Report {
  id: string;
  text: string;
  image: string;
  location: string;
  status: ReportStatus;
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

interface PaginationData {
  data: Report[];
  current_page: number;
  last_page: number;
  per_page: number;
  total: number;
}

interface CountData {
  total: number;
  unverified: number;
  verified: number;
  rejected: number;
  handled: number;
  completed: number;
}

interface CountApiResponse {
  status: boolean;
  message: string;
  data: CountData;
}

export default function AdminDashboard() {
  const router = useRouter();
  const [reports, setReports] = useState<Report[]>([]);
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

  // Pagination states
  const [currentPage, setCurrentPage] = useState(1);
  const [perPage, setPerPage] = useState(10);
  const [totalPages, setTotalPages] = useState(1);
  const [totalRecords, setTotalRecords] = useState(0);
  const [search, setSearch] = useState('');
  const [searchInput, setSearchInput] = useState('');

  // Count data state
  const [countData, setCountData] = useState<CountData>({
    total: 0,
    unverified: 0,
    verified: 0,
    rejected: 0,
    handled: 0,
    completed: 0
  });
  const [countLoading, setCountLoading] = useState(true);

  // Fetch count data from API
  const fetchCountData = async () => {
    try {
      setCountLoading(true);
      const token = getFromSessionStorage('access_token');
      const headers = token ? { Authorization: `Bearer ${token}` } : {};

      console.log('API Request (Count):', `${backendUrl}/api/reports/count`);

      const response = await axios.get<CountApiResponse>(
        `${backendUrl}/api/reports/count`,
        { headers }
      );

      if (response.status === 200 && response.data.status) {
        setCountData(response.data.data);
        // Also update totalRecords from count data for consistency
        setTotalRecords(response.data.data.total);
      } else {
        console.error('Failed to fetch count data:', response.data.message);
      }
    } catch (error) {
      console.error('Error fetching count data:', error);
      // Check if it's an auth error and handle redirect
      if (!handleAuthError(error)) {
        // If not auth error, just log it (don't show error for count data)
        console.warn('Count data fetch failed, using fallback');
      }
    } finally {
      setCountLoading(false);
    }
  };

  // Fetch reports from API with pagination
  const fetchReports = async (page = 1, searchTerm = '', itemsPerPage = 10) => {
    try {
      setLoading(true);
      setError(null);
      const token = getFromSessionStorage('access_token');
      const headers = token ? { Authorization: `Bearer ${token}` } : {};

      const params = new URLSearchParams({
        page: page.toString(),
        per_page: itemsPerPage.toString(),
        ...(searchTerm && { search: searchTerm })
      });

      console.log('API Request:', `${backendUrl}/api/reports?${params}`);

      const response = await axios.get(
        `${backendUrl}/api/reports?${params}`,
        { headers }
      );

      console.log('API Response:', response.data);

      if (response.status === 200 && response.data.status) {
        // Handle different possible response structures
        let reportsData, paginationMeta;

        if (response.data.data && response.data.data.data) {
          // Structure: { status, message, data: { data: [...], current_page, last_page, etc } }
          const { data, current_page, last_page, per_page, total } = response.data.data;
          reportsData = data;
          paginationMeta = { current_page, last_page, per_page, total };
        } else if (response.data.data && Array.isArray(response.data.data)) {
          // Structure: { status, message, data: [...] } (no pagination)
          reportsData = response.data.data;
          paginationMeta = { current_page: 1, last_page: 1, per_page: reportsData.length, total: reportsData.length };
        } else {
          throw new Error('Unexpected API response structure');
        }

        setReports(reportsData);
        setTotalPages(paginationMeta.last_page);
        setTotalRecords(paginationMeta.total);

        // Verify that the response matches our request
        if (paginationMeta.per_page !== itemsPerPage) {
          setPerPage(paginationMeta.per_page);
        }
      } else {
        throw new Error('Failed to fetch reports');
      }
    } catch (error) {
      console.error('Error fetching reports:', error);
      // Check if it's an auth error and handle redirect
      if (!handleAuthError(error)) {
        // If not auth error, show general error message
        setError('Failed to load reports. Please try again later.');
        setReports([]);
      }
    } finally {
      setLoading(false);
    }
  };

  useEffect(() => {
    console.log('Fetching reports:', { currentPage, search, perPage });
    fetchReports(currentPage, search, perPage);
  }, [currentPage, search, perPage]);

  // Fetch count data on component mount and when status changes
  useEffect(() => {
    fetchCountData();
  }, []);

  // Refetch count data after successful status update
  const refetchCountAfterStatusUpdate = () => {
    fetchCountData();
  };

  // Debounced search effect
  useEffect(() => {
    const debounceTimer = setTimeout(() => {
      if (searchInput !== search) {
        setSearch(searchInput);
      }
    }, 500); // 500ms delay

    return () => clearTimeout(debounceTimer);
  }, [searchInput, search]);

  const getStatusBadge = (status: ReportStatus) => {
    switch (status) {
      case 'verified':
        return <Badge variant="verified">Verified</Badge>;
      case 'unverified':
        return <Badge variant="unverified">Unverified</Badge>;
      case 'handled':
        return <Badge variant="handled">Handled</Badge>;
      case 'rejected':
        return <Badge variant="rejected">Rejected</Badge>;
      default:
        return <Badge variant="secondary">{status}</Badge>;
    }
  };

  // Add loading state for status updates
  const [updatingStatus, setUpdatingStatus] = useState<string | null>(null);

  const handleStatusChange = async (reportId: string, newStatus: ReportStatus) => {
    try {
      setUpdatingStatus(reportId);
      const token = getFromSessionStorage('access_token');
      const headers = token ? { Authorization: `Bearer ${token}` } : {};

      // Send API request to update status
      const response = await axios.post(
        `${backendUrl}/api/reports/${reportId}/status`,
        { status: newStatus },
        { headers }
      );

      if (response.status === 200) {
        // Update local state only if API call succeeds
        setReports(prev =>
          prev.map(report =>
            report.id === reportId
              ? { ...report, status: newStatus }
              : report
          )
        );
        console.log(`Status updated successfully for report ${reportId} to ${newStatus}`);

        // Refetch count data to get updated statistics
        refetchCountAfterStatusUpdate();
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
      setUpdatingStatus(null);
    }
  };

  const navigateToReportDetail = (reportId: string) => {
    router.push(`/report/${reportId}`);
  };

  const handlePageChange = (page: number) => {
    if (page >= 1 && page <= totalPages && page !== currentPage) {
      console.log('Page change:', { from: currentPage, to: page });
      setCurrentPage(page);
    }
  };

  const handlePerPageChange = (newPerPage: number) => {
    console.log('Per page change:', { from: perPage, to: newPerPage });
    setPerPage(newPerPage);
    setCurrentPage(1); // Reset to first page when changing per page
  };

  const handleSearchChange = (searchTerm: string) => {
    console.log('Search change:', { from: searchInput, to: searchTerm });
    setSearchInput(searchTerm);
    if (currentPage !== 1) {
      setCurrentPage(1); // Reset to first page when searching
    }
  };

  const getPaginationNumbers = () => {
    const numbers = [];
    const maxVisible = 5;

    if (totalPages <= maxVisible) {
      // If total pages is less than max visible, show all pages
      for (let i = 1; i <= totalPages; i++) {
        numbers.push(i);
      }
    } else {
      // Calculate start and end based on current page
      let start = Math.max(1, currentPage - Math.floor(maxVisible / 2));
      const end = Math.min(totalPages, start + maxVisible - 1);

      // Adjust start if we're near the end
      if (end - start + 1 < maxVisible) {
        start = Math.max(1, end - maxVisible + 1);
      }

      for (let i = start; i <= end; i++) {
        numbers.push(i);
      }
    }

    return numbers;
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

  if (error) {
    return (
      <Layout>
        <Seo templateTitle="Admin Dashboard - RAPID" />
        <main className="min-h-screen bg-gray-50 p-6">
          <div className="mx-auto max-w-7xl">
            <Card className="bg-red-50 border-red-200">
              <CardContent className="p-6 text-center">
                <div className="text-red-600 text-xl mb-2">⚠️ Error</div>
                <p className="text-red-800">{error}</p>
                <Button
                  onClick={() => window.location.reload()}
                  className="mt-4"
                  variant="outline"
                >
                  Retry
                </Button>
              </CardContent>
            </Card>
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
            <div className="flex flex-col lg:flex-row lg:items-center lg:justify-between gap-4">
              <div>
                <h1 className="text-3xl font-bold text-gray-900">RAPID Admin Dashboard</h1>
                <p className="text-gray-600 mt-2">Manage and monitor user reports</p>
              </div>
              <div className="flex flex-col sm:flex-row items-start sm:items-center gap-4">
                <div className="flex items-center space-x-2">
                  <span className="text-sm text-gray-500">Total:</span>
                  <Badge variant="info">{countData.total}</Badge>
                </div>
                <div className="w-full sm:w-80">
                  <Input
                    placeholder="Search reports..."
                    value={searchInput}
                    onChange={(e) => handleSearchChange(e.target.value)}
                    className="w-full"
                  />
                </div>
              </div>
            </div>
          </div>

          {/* Stats Cards */}
          <div className="grid grid-cols-1 md:grid-cols-5 gap-6 mb-8">
            <Card>
              <CardHeader className="pb-2">
                <CardTitle className="text-sm font-medium text-gray-600">Total Reports</CardTitle>
              </CardHeader>
              <CardContent>
                <div className="text-2xl font-bold text-gray-900">
                  {countLoading ? (
                    <div className="animate-pulse bg-gray-300 h-8 w-12 rounded"></div>
                  ) : (
                    countData.total
                  )}
                </div>
              </CardContent>
            </Card>
            <Card>
              <CardHeader className="pb-2">
                <CardTitle className="text-sm font-medium text-gray-600">Unverified</CardTitle>
              </CardHeader>
              <CardContent>
                <div className="text-2xl font-bold text-yellow-600">
                  {countLoading ? (
                    <div className="animate-pulse bg-gray-300 h-8 w-12 rounded"></div>
                  ) : (
                    countData.unverified
                  )}
                </div>
              </CardContent>
            </Card>
            <Card>
              <CardHeader className="pb-2">
                <CardTitle className="text-sm font-medium text-gray-600">Verified</CardTitle>
              </CardHeader>
              <CardContent>
                <div className="text-2xl font-bold text-green-600">
                  {countLoading ? (
                    <div className="animate-pulse bg-gray-300 h-8 w-12 rounded"></div>
                  ) : (
                    countData.verified
                  )}
                </div>
              </CardContent>
            </Card>
            <Card>
              <CardHeader className="pb-2">
                <CardTitle className="text-sm font-medium text-gray-600">Handled</CardTitle>
              </CardHeader>
              <CardContent>
                <div className="text-2xl font-bold text-blue-600">
                  {countLoading ? (
                    <div className="animate-pulse bg-gray-300 h-8 w-12 rounded"></div>
                  ) : (
                    countData.handled
                  )}
                </div>
              </CardContent>
            </Card>
            <Card>
              <CardHeader className="pb-2">
                <CardTitle className="text-sm font-medium text-gray-600">Rejected</CardTitle>
              </CardHeader>
              <CardContent>
                <div className="text-2xl font-bold text-red-600">
                  {countLoading ? (
                    <div className="animate-pulse bg-gray-300 h-8 w-12 rounded"></div>
                  ) : (
                    countData.rejected
                  )}
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
                      <TableHead className="min-w-[150px]">Location</TableHead>
                      <TableHead>User</TableHead>
                      <TableHead>Tag</TableHead>
                      <TableHead>Created</TableHead>
                      <TableHead className="w-[200px]">Actions</TableHead>
                    </TableRow>
                  </TableHeader>
                  <TableBody>
                    {reports.map((report) => (
                      <TableRow
                        key={report.id}
                        onClick={() => navigateToReportDetail(report.id)}
                        className="cursor-pointer hover:bg-gray-50"
                      >
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
                          {report.pred_confidence > 0 ? (
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
                            <span className="text-gray-400">0%</span>
                          )}
                        </TableCell>
                        <TableCell>
                          <span className="text-sm text-gray-600">{report.location}</span>
                        </TableCell>
                        <TableCell>
                          <div>
                            <p className="text-sm font-medium">{report.username}</p>
                            <p className="text-xs text-gray-500">ID: {report.user_id}</p>
                          </div>
                        </TableCell>
                        <TableCell>
                          {report.tag && report.tag.class ? (
                            <div className='flex flex-wrap gap-1'>
                              {report.tag.class.split(',').map((tag, index) => (
                                <Badge key={index} variant="outline">
                                  {tag.trim()}
                                </Badge>
                              ))}
                            </div>
                          ) : (
                            <span className="text-gray-400">No Tag</span>
                          )}
                        </TableCell>
                        <TableCell>
                          <span className="text-xs text-gray-500">
                            {report.created_at ?
                              new Date(report.created_at).toLocaleDateString() :
                              'No date'
                            }
                          </span>
                        </TableCell>
                        <TableCell>
                          <div className="flex space-x-1">
                            <Button
                              size="sm"
                              variant="outline"
                              className="text-green-600 border-green-600 hover:bg-green-50"
                              onClick={async (e) => {
                                e.stopPropagation();
                                await handleStatusChange(report.id, 'verified');
                              }}
                              disabled={report.status === 'verified' || updatingStatus === report.id}
                              title="Mark as Verified"
                            >
                              {updatingStatus === report.id ? '⏳' : '✓'}
                            </Button>
                            <Button
                              size="sm"
                              variant="outline"
                              className="text-blue-600 border-blue-600 hover:bg-blue-50"
                              onClick={async (e) => {
                                e.stopPropagation();
                                await handleStatusChange(report.id, 'handled');
                              }}
                              disabled={report.status === 'handled' || updatingStatus === report.id}
                              title="Mark as Handled"
                            >
                              {updatingStatus === report.id ? '⏳' : '✉'}
                            </Button>
                            <Button
                              size="sm"
                              variant="outline"
                              className="text-red-600 border-red-600 hover:bg-red-50"
                              onClick={async (e) => {
                                e.stopPropagation();
                                await handleStatusChange(report.id, 'rejected');
                              }}
                              disabled={report.status === 'rejected' || updatingStatus === report.id}
                              title="Mark as Rejected"
                            >
                              {updatingStatus === report.id ? '⏳' : '✕'}
                            </Button>
                            <Button
                              size="sm"
                              variant="outline"
                              className="text-yellow-600 border-yellow-600 hover:bg-yellow-50"
                              onClick={async (e) => {
                                e.stopPropagation();
                                await handleStatusChange(report.id, 'unverified');
                              }}
                              disabled={report.status === 'unverified' || updatingStatus === report.id}
                              title="Mark as Unverified"
                            >
                              {updatingStatus === report.id ? '⏳' : '?'}
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

          {/* Pagination Controls */}
          <Card className="mt-6">
            <CardContent className="p-6">
              <div className="flex flex-col sm:flex-row items-center justify-between gap-4">
                {/* Items per page selector */}
                <div className="flex items-center space-x-2">
                  <span className="text-sm text-gray-700">Show:</span>
                  <select
                    value={perPage}
                    onChange={(e) => handlePerPageChange(Number(e.target.value))}
                    className="border border-gray-300 rounded-md px-4 py-2 text-base focus:outline-none focus:ring-2 focus:ring-blue-500 min-w-[70px]"
                    style={{ fontSize: '1rem', height: '2.5rem' }}
                  >
                    <option value={5}>5</option>
                    <option value={10}>10</option>
                    <option value={25}>25</option>
                    <option value={50}>50</option>
                  </select>
                  <span className="text-sm text-gray-700">per page</span>
                </div>

                {/* Page info */}
                <div className="text-sm text-gray-700">
                  {totalRecords > 0 ? (
                    <>Showing {((currentPage - 1) * perPage) + 1} to {Math.min(currentPage * perPage, totalRecords)} of {totalRecords} results</>
                  ) : (
                    'No results found'
                  )}
                </div>

                {/* Pagination buttons - only show if multiple pages */}
                {totalPages > 1 && (
                  <div className="flex items-center space-x-1">
                    <Button
                      variant="outline"
                      size="sm"
                      onClick={() => handlePageChange(1)}
                      disabled={currentPage === 1}
                      className="px-2"
                    >
                      ««
                    </Button>
                    <Button
                      variant="outline"
                      size="sm"
                      onClick={() => handlePageChange(currentPage - 1)}
                      disabled={currentPage === 1}
                      className="px-2"
                    >
                      «
                    </Button>

                    {getPaginationNumbers().map((page) => (
                      <Button
                        key={page}
                        variant={currentPage === page ? "default" : "outline"}
                        size="sm"
                        onClick={() => handlePageChange(page)}
                        className="px-3"
                      >
                        {page}
                      </Button>
                    ))}

                    <Button
                      variant="outline"
                      size="sm"
                      onClick={() => handlePageChange(currentPage + 1)}
                      disabled={currentPage === totalPages}
                      className="px-2"
                    >
                      »
                    </Button>
                    <Button
                      variant="outline"
                      size="sm"
                      onClick={() => handlePageChange(totalPages)}
                      disabled={currentPage === totalPages}
                      className="px-2"
                    >
                      »»
                    </Button>
                  </div>
                )}
              </div>
            </CardContent>
          </Card>
        </div>
      </main>
    </Layout>
  );
}
