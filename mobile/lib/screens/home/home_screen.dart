import 'package:flutter/material.dart';
import '../../utils/constants.dart';
import '../report/create_report_screen.dart';
import '../auth/login_screen.dart';
import '../../services/auth_service.dart';
import '../../services/report_service.dart';

class HomeScreen extends StatefulWidget {
  final String token;

  const HomeScreen({super.key, required this.token});

  @override
  State<HomeScreen> createState() => _HomeScreenState();
}

class _HomeScreenState extends State<HomeScreen> with TickerProviderStateMixin {
  late TabController _tabController;
  int _selectedIndex = 0;
  final ReportService _reportService = ReportService();

  // Data state
  List<Map<String, dynamic>> _allReports = [];
  List<Map<String, dynamic>> _userReports = [];
  bool _isLoading = false;
  bool _isLoadingMore = false;
  String? _errorMessage;

  // Pagination state
  int _currentPage = 1;
  int _userCurrentPage = 1;
  bool _hasMoreData = true;
  bool _hasMoreUserData = true;
  final int _perPage = 10;

  // Search state
  String _searchQuery = '';
  bool _isSearching = false;
  final TextEditingController _searchController = TextEditingController();

  // Scroll controllers for infinite scroll
  final ScrollController _scrollController = ScrollController();
  final ScrollController _userScrollController = ScrollController();
  @override
  void initState() {
    super.initState();
    _tabController = TabController(length: 3, vsync: this);
    _loadReports();

    // Setup scroll listeners for infinite scroll
    _scrollController.addListener(_onScroll);
    _userScrollController.addListener(_onUserScroll);
  }

  // Scroll listener for infinite scroll
  void _onScroll() {
    if (_scrollController.position.pixels ==
        _scrollController.position.maxScrollExtent) {
      if (!_isLoadingMore && _hasMoreData) {
        _loadMoreReports();
      }
    }
  }

  void _onUserScroll() {
    if (_userScrollController.position.pixels ==
        _userScrollController.position.maxScrollExtent) {
      if (!_isLoadingMore && _hasMoreUserData) {
        _loadMoreUserReports();
      }
    }
  }

  // Load reports from API (first page or refresh)
  Future<void> _loadReports({bool isRefresh = false}) async {
    if (isRefresh) {
      _currentPage = 1;
      _hasMoreData = true;
    }

    setState(() {
      _isLoading = true;
      _errorMessage = null;
      if (isRefresh) _allReports.clear();
    });

    try {
      final response = await _reportService.getAllReports(
        page: _currentPage,
        perPage: _perPage,
        search: _searchQuery.isNotEmpty ? _searchQuery : null,
      );

      if (response.isSuccess && response.data != null) {
        // Parse nested response structure: response.data.data
        final responseDataMap = response.data!;
        final dataWrapper = responseDataMap['data'] as Map<String, dynamic>?;
        final List<dynamic> reportsData = dataWrapper?['data'] ?? [];

        final List<Map<String, dynamic>> reports = reportsData
            .map((report) => report as Map<String, dynamic>)
            .toList();

        setState(() {
          if (isRefresh) {
            _allReports = reports;
          } else {
            _allReports.addAll(reports);
          }
          _isLoading = false;

          // Check if there's more data
          _hasMoreData = reports.length == _perPage;
          if (!isRefresh) _currentPage++;
        });
      } else {
        setState(() {
          _isLoading = false;
          _errorMessage = response.message;
        });
      }
    } catch (e) {
      setState(() {
        _isLoading = false;
        _errorMessage = 'Failed to load reports: ${e.toString()}';
      });
    }
  }

  // Load more reports for infinite scroll
  Future<void> _loadMoreReports() async {
    if (_isLoadingMore || !_hasMoreData) return;

    setState(() => _isLoadingMore = true);

    try {
      final response = await _reportService.getAllReports(
        page: _currentPage,
        perPage: _perPage,
        search: _searchQuery.isNotEmpty ? _searchQuery : null,
      );

      if (response.isSuccess && response.data != null) {
        // Parse nested response structure: response.data.data
        final responseDataMap = response.data!;
        final dataWrapper = responseDataMap['data'] as Map<String, dynamic>?;
        final List<dynamic> reportsData = dataWrapper?['data'] ?? [];

        final List<Map<String, dynamic>> reports = reportsData
            .map((report) => report as Map<String, dynamic>)
            .toList();

        setState(() {
          _allReports.addAll(reports);
          _currentPage++;
          _hasMoreData = reports.length == _perPage;
          _isLoadingMore = false;
        });
      } else {
        setState(() => _isLoadingMore = false);
      }
    } catch (e) {
      setState(() => _isLoadingMore = false);
    }
  }

  // Load user reports
  Future<void> _loadUserReports({bool isRefresh = false}) async {
    if (isRefresh) {
      _userCurrentPage = 1;
      _hasMoreUserData = true;
    }

    setState(() {
      _isLoading = true;
      _errorMessage = null;
      if (isRefresh) _userReports.clear();
    });

    try {
      final response = await _reportService.getUserReports(
        page: _userCurrentPage,
        perPage: _perPage,
        search: _searchQuery.isNotEmpty ? _searchQuery : null,
      );

      if (response.isSuccess && response.data != null) {
        // Parse nested response structure: response.data.data
        final responseDataMap = response.data!;
        final dataWrapper = responseDataMap['data'] as Map<String, dynamic>?;
        final List<dynamic> reportsData = dataWrapper?['data'] ?? [];

        final List<Map<String, dynamic>> reports = reportsData
            .map((report) => report as Map<String, dynamic>)
            .toList();

        setState(() {
          if (isRefresh) {
            _userReports = reports;
          } else {
            _userReports.addAll(reports);
          }
          _isLoading = false;

          // Check if there's more data
          _hasMoreUserData = reports.length == _perPage;
          if (!isRefresh) _userCurrentPage++;
        });
      } else {
        setState(() {
          _isLoading = false;
          _errorMessage = response.message;
        });
      }
    } catch (e) {
      setState(() {
        _isLoading = false;
        _errorMessage = 'Failed to load user reports: ${e.toString()}';
      });
    }
  }

  // Load more user reports for infinite scroll
  Future<void> _loadMoreUserReports() async {
    if (_isLoadingMore || !_hasMoreUserData) return;

    setState(() => _isLoadingMore = true);

    try {
      final response = await _reportService.getUserReports(
        page: _userCurrentPage,
        perPage: _perPage,
        search: _searchQuery.isNotEmpty ? _searchQuery : null,
      );

      if (response.isSuccess && response.data != null) {
        // Parse nested response structure: response.data.data
        final responseDataMap = response.data!;
        final dataWrapper = responseDataMap['data'] as Map<String, dynamic>?;
        final List<dynamic> reportsData = dataWrapper?['data'] ?? [];

        final List<Map<String, dynamic>> reports = reportsData
            .map((report) => report as Map<String, dynamic>)
            .toList();

        setState(() {
          _userReports.addAll(reports);
          _userCurrentPage++;
          _hasMoreUserData = reports.length == _perPage;
          _isLoadingMore = false;
        });
      } else {
        setState(() => _isLoadingMore = false);
      }
    } catch (e) {
      setState(() => _isLoadingMore = false);
    }
  }

  // Refresh reports
  Future<void> _refreshReports() async {
    await _loadReports(isRefresh: true);
  }

  // Search functionality
  void _onSearchChanged(String query) {
    setState(() {
      _searchQuery = query;
      _isSearching = true;
    });

    // Debounce search
    Future.delayed(const Duration(milliseconds: 500), () {
      if (_searchQuery == query) {
        _loadReports(isRefresh: true).then((_) {
          if (mounted) {
            setState(() => _isSearching = false);
          }
        });
      }
    });
  }

  void _clearSearch() {
    _searchController.clear();
    setState(() {
      _searchQuery = '';
      _isSearching = false;
    });
    _loadReports(isRefresh: true);
  }

  @override
  Widget build(BuildContext context) {
    return Scaffold(
      backgroundColor: AppColors.background,
      body: NestedScrollView(
        headerSliverBuilder: (context, innerBoxIsScrolled) {
          return [
            SliverAppBar(
              expandedHeight: 120,
              floating: false,
              pinned: true,
              backgroundColor: AppColors.primary,
              foregroundColor: Colors.white,
              elevation: 0,
              flexibleSpace: FlexibleSpaceBar(
                // title: const Text(
                // 'Reports',
                // style: TextStyle(fontWeight: FontWeight.bold, fontSize: 20),
                // ),
                background: Container(
                  decoration: BoxDecoration(
                    gradient: LinearGradient(
                      begin: Alignment.topLeft,
                      end: Alignment.bottomRight,
                      colors: [
                        AppColors.primary,
                        AppColors.primary.withValues(alpha: 0.8),
                      ],
                    ),
                  ),
                ),
              ),
              actions: [
                IconButton(
                  onPressed: () {
                    _showUserMenu(context);
                  },
                  icon: const CircleAvatar(
                    radius: 16,
                    backgroundColor: Colors.white,
                    child: Icon(
                      Icons.person,
                      color: AppColors.primary,
                      size: 20,
                    ),
                  ),
                ),
              ],
              bottom: TabBar(
                controller: _tabController,
                indicatorColor: Colors.white,
                indicatorWeight: 3,
                labelColor: Colors.white,
                unselectedLabelColor: Colors.white70,
                labelStyle: const TextStyle(fontWeight: FontWeight.w600),
                tabs: const [
                  Tab(text: 'All'),
                  Tab(text: 'Trending'),
                  Tab(text: 'Your Report'),
                ],
              ),
            ),
          ];
        },
        body: Column(
          children: [
            // Search Bar
            Container(
              padding: const EdgeInsets.all(16),
              color: Colors.white,
              child: TextField(
                controller: _searchController,
                onChanged: _onSearchChanged,
                decoration: InputDecoration(
                  hintText: 'Search reports...',
                  prefixIcon: const Icon(Icons.search, color: AppColors.accent),
                  suffixIcon: _isSearching
                      ? const SizedBox(
                          width: 20,
                          height: 20,
                          child: Padding(
                            padding: EdgeInsets.all(12),
                            child: CircularProgressIndicator(strokeWidth: 2),
                          ),
                        )
                      : _searchQuery.isNotEmpty
                      ? IconButton(
                          icon: const Icon(
                            Icons.clear,
                            color: AppColors.accent,
                          ),
                          onPressed: _clearSearch,
                        )
                      : null,
                  border: OutlineInputBorder(
                    borderRadius: BorderRadius.circular(25),
                    borderSide: const BorderSide(color: AppColors.border),
                  ),
                  enabledBorder: OutlineInputBorder(
                    borderRadius: BorderRadius.circular(25),
                    borderSide: const BorderSide(color: AppColors.border),
                  ),
                  focusedBorder: OutlineInputBorder(
                    borderRadius: BorderRadius.circular(25),
                    borderSide: const BorderSide(
                      color: AppColors.primary,
                      width: 2,
                    ),
                  ),
                  filled: true,
                  fillColor: AppColors.background,
                  contentPadding: const EdgeInsets.symmetric(
                    horizontal: 20,
                    vertical: 16,
                  ),
                ),
              ),
            ),
            // Tab Content
            Expanded(
              child: TabBarView(
                controller: _tabController,
                children: [
                  _buildReportsList(_allReports, _scrollController),
                  _buildEmptyState('Trending reports will appear here'),
                  _buildUserReportsList(),
                ],
              ),
            ),
          ],
        ),
      ),
      floatingActionButton: FloatingActionButton.extended(
        onPressed: () async {
          final result = await Navigator.push(
            context,
            MaterialPageRoute(
              builder: (context) => CreateReportScreen(token: widget.token),
            ),
          );

          // Refresh reports if a new report was created
          if (result == true) {
            _loadReports();
          }
        },
        backgroundColor: AppColors.primary,
        foregroundColor: Colors.white,
        icon: const Icon(Icons.add),
        label: const Text('Report'),
        elevation: 4,
      ),
      bottomNavigationBar: BottomNavigationBar(
        currentIndex: _selectedIndex,
        onTap: (index) => setState(() => _selectedIndex = index),
        type: BottomNavigationBarType.fixed,
        backgroundColor: Colors.white,
        selectedItemColor: AppColors.primary,
        unselectedItemColor: AppColors.accent,
        elevation: 8,
        items: const [
          BottomNavigationBarItem(
            icon: Icon(Icons.home_outlined),
            activeIcon: Icon(Icons.home),
            label: 'Home',
          ),
          BottomNavigationBarItem(
            icon: Icon(Icons.search_outlined),
            activeIcon: Icon(Icons.search),
            label: 'Search',
          ),
          BottomNavigationBarItem(
            icon: Icon(Icons.notifications_outlined),
            activeIcon: Icon(Icons.notifications),
            label: 'Notifications',
          ),
          BottomNavigationBarItem(
            icon: Icon(Icons.person_outline),
            activeIcon: Icon(Icons.person),
            label: 'Profile',
          ),
        ],
      ),
    );
  }

  Widget _buildReportsList(
    List<Map<String, dynamic>> reports,
    ScrollController scrollController,
  ) {
    if (_isLoading && reports.isEmpty) {
      return const Center(child: CircularProgressIndicator());
    }

    if (_errorMessage != null && reports.isEmpty) {
      return Center(
        child: Column(
          mainAxisAlignment: MainAxisAlignment.center,
          children: [
            Icon(
              Icons.error_outline,
              size: 80,
              color: AppColors.error.withValues(alpha: 0.5),
            ),
            const SizedBox(height: 16),
            Text(
              _errorMessage!,
              style: AppTextStyles.body2.copyWith(color: AppColors.error),
              textAlign: TextAlign.center,
            ),
            const SizedBox(height: 16),
            ElevatedButton(
              onPressed: () => _loadReports(isRefresh: true),
              child: const Text('Retry'),
            ),
          ],
        ),
      );
    }

    if (reports.isEmpty) {
      return _buildEmptyState('No reports available');
    }

    return RefreshIndicator(
      onRefresh: _refreshReports,
      child: ListView.builder(
        controller: scrollController,
        padding: const EdgeInsets.symmetric(vertical: 8),
        itemCount:
            reports.length +
            (_shouldShowLoadingMore(reports, scrollController) ? 1 : 0),
        itemBuilder: (context, index) {
          if (index == reports.length) {
            // Loading indicator for infinite scroll
            return Container(
              padding: const EdgeInsets.all(16),
              alignment: Alignment.center,
              child: _isLoadingMore
                  ? const CircularProgressIndicator()
                  : const SizedBox.shrink(),
            );
          }

          final report = reports[index];
          return _ReportCard(
            report: report,
            onLike: () =>
                _toggleLike(report['id']?.toString() ?? '', index, reports),
            onShare: () => _shareReport(report),
          );
        },
      ),
    );
  }

  bool _shouldShowLoadingMore(
    List<Map<String, dynamic>> reports,
    ScrollController scrollController,
  ) {
    if (scrollController == _scrollController) {
      return _hasMoreData && !_isLoading;
    } else {
      return _hasMoreUserData && !_isLoading;
    }
  }

  Widget _buildUserReportsList() {
    return FutureBuilder<void>(
      future: _userReports.isEmpty ? _loadUserReports() : null,
      builder: (context, snapshot) {
        if (snapshot.connectionState == ConnectionState.waiting &&
            _userReports.isEmpty) {
          return const Center(child: CircularProgressIndicator());
        }

        return _buildReportsList(_userReports, _userScrollController);
      },
    );
  }

  Widget _buildEmptyState(String message) {
    return Center(
      child: Column(
        mainAxisAlignment: MainAxisAlignment.center,
        children: [
          Icon(
            Icons.inbox_outlined,
            size: 80,
            color: AppColors.accent.withValues(alpha: 0.5),
          ),
          const SizedBox(height: 16),
          Text(
            message,
            style: AppTextStyles.body2.copyWith(
              color: AppColors.accent.withValues(alpha: 0.7),
            ),
          ),
        ],
      ),
    );
  }

  Future<void> _toggleLike(
    String reportId,
    int index,
    List<Map<String, dynamic>> reports,
  ) async {
    try {
      // Optimistic update - update UI immediately
      setState(() {
        reports[index]['isLiked'] = !reports[index]['isLiked'];
        if (reports[index]['isLiked']) {
          reports[index]['likes']++;
        } else {
          reports[index]['likes']--;
        }
      });

      // Call API
      final response = await _reportService.toggleLike(reportId);

      if (!response.isSuccess) {
        // Revert if API call fails
        setState(() {
          reports[index]['isLiked'] = !reports[index]['isLiked'];
          if (reports[index]['isLiked']) {
            reports[index]['likes']++;
          } else {
            reports[index]['likes']--;
          }
        });

        if (mounted) {
          ScaffoldMessenger.of(context).showSnackBar(
            SnackBar(
              content: Text('Failed to like report: ${response.message}'),
              backgroundColor: AppColors.error,
            ),
          );
        }
      }
    } catch (e) {
      // Revert on error
      setState(() {
        reports[index]['isLiked'] = !reports[index]['isLiked'];
        if (reports[index]['isLiked']) {
          reports[index]['likes']++;
        } else {
          reports[index]['likes']--;
        }
      });

      if (mounted) {
        ScaffoldMessenger.of(context).showSnackBar(
          SnackBar(
            content: Text('Error: ${e.toString()}'),
            backgroundColor: AppColors.error,
          ),
        );
      }
    }
  }

  void _shareReport(Map<String, dynamic> report) {
    ScaffoldMessenger.of(context).showSnackBar(
      const SnackBar(
        content: Text('Share feature coming soon!'),
        backgroundColor: AppColors.success,
      ),
    );
  }

  void _showUserMenu(BuildContext context) {
    showModalBottomSheet(
      context: context,
      shape: const RoundedRectangleBorder(
        borderRadius: BorderRadius.vertical(top: Radius.circular(20)),
      ),
      builder: (context) => Container(
        padding: const EdgeInsets.all(16),
        child: Column(
          mainAxisSize: MainAxisSize.min,
          children: [
            ListTile(
              leading: const Icon(Icons.person, color: AppColors.primary),
              title: const Text('Profile'),
              onTap: () {
                Navigator.pop(context);
                ScaffoldMessenger.of(context).showSnackBar(
                  const SnackBar(content: Text('Profile feature coming soon!')),
                );
              },
            ),
            ListTile(
              leading: const Icon(Icons.settings, color: AppColors.accent),
              title: const Text('Settings'),
              onTap: () {
                Navigator.pop(context);
                ScaffoldMessenger.of(context).showSnackBar(
                  const SnackBar(
                    content: Text('Settings feature coming soon!'),
                  ),
                );
              },
            ),
            ListTile(
              leading: const Icon(Icons.logout, color: AppColors.error),
              title: const Text('Logout'),
              onTap: () {
                Navigator.pop(context);
                _logout();
              },
            ),
          ],
        ),
      ),
    );
  }

  void _logout() {
    showDialog(
      context: context,
      builder: (context) => AlertDialog(
        title: const Text('Logout'),
        content: const Text('Are you sure you want to logout?'),
        actions: [
          TextButton(
            onPressed: () => Navigator.pop(context),
            child: const Text('Cancel'),
          ),
          ElevatedButton(
            onPressed: () async {
              Navigator.pop(context);
              final authService = AuthService();
              await authService.logout();
              Navigator.pushAndRemoveUntil(
                context,
                MaterialPageRoute(builder: (context) => const LoginScreen()),
                (route) => false,
              );
            },
            style: ElevatedButton.styleFrom(backgroundColor: AppColors.error),
            child: const Text('Logout'),
          ),
        ],
      ),
    );
  }

  @override
  void dispose() {
    _tabController.dispose();
    _searchController.dispose();
    _scrollController.dispose();
    _userScrollController.dispose();
    super.dispose();
  }
}

class _ReportCard extends StatelessWidget {
  final Map<String, dynamic> report;
  final VoidCallback onLike;
  final VoidCallback onShare;

  const _ReportCard({
    required this.report,
    required this.onLike,
    required this.onShare,
  });

  @override
  Widget build(BuildContext context) {
    // Handle different data formats from API
    final String userName =
        report['user']?['name'] ??
        report['username'] ??
        report['author']?['name'] ??
        'Anonymous';
    final String userLocation =
        report['location'] ?? report['user']?['location'] ?? 'Unknown Location';
    final String timeAgo =
        report['time'] ??
        report['created_at'] ??
        report['time_ago'] ??
        'Unknown time';
    final String reportText =
        report['text'] ?? report['content'] ?? report['description'] ?? '';
    final int likesCount = report['likes'] ?? report['likes_count'] ?? 0;
    final int sharesCount = report['shares'] ?? report['shares_count'] ?? 0;
    final bool isLiked =
        report['isLiked'] ??
        report['is_liked'] ??
        report['user_liked'] ??
        false;
    final String? imageUrl =
        report['image'] ?? report['image_url'] ?? report['photo_url'];

    return Container(
      margin: const EdgeInsets.symmetric(horizontal: 12, vertical: 6),
      decoration: BoxDecoration(
        color: AppColors.cardBackground,
        borderRadius: BorderRadius.circular(16),
        boxShadow: [
          BoxShadow(
            color: Colors.black.withValues(alpha: 0.05),
            blurRadius: 10,
            offset: const Offset(0, 2),
          ),
        ],
      ),
      child: Column(
        crossAxisAlignment: CrossAxisAlignment.start,
        children: [
          // Header
          Padding(
            padding: const EdgeInsets.all(16),
            child: Row(
              children: [
                CircleAvatar(
                  radius: 22,
                  backgroundColor: AppColors.primary.withValues(alpha: 0.1),
                  child: Text(
                    userName.isNotEmpty ? userName[0].toUpperCase() : 'A',
                    style: const TextStyle(
                      color: AppColors.primary,
                      fontWeight: FontWeight.bold,
                    ),
                  ),
                ),
                const SizedBox(width: 12),
                Expanded(
                  child: Column(
                    crossAxisAlignment: CrossAxisAlignment.start,
                    children: [
                      Text(
                        userName,
                        style: AppTextStyles.body1.copyWith(
                          fontWeight: FontWeight.w600,
                        ),
                      ),
                      Row(
                        children: [
                          Icon(
                            Icons.location_on,
                            size: 14,
                            color: AppColors.accent,
                          ),
                          const SizedBox(width: 4),
                          Flexible(
                            child: Text(
                              userLocation,
                              style: AppTextStyles.caption,
                              overflow: TextOverflow.ellipsis,
                            ),
                          ),
                          const Text(
                            ' • ',
                            style: TextStyle(color: AppColors.accent),
                          ),
                          Text(timeAgo, style: AppTextStyles.caption),
                        ],
                      ),
                    ],
                  ),
                ),
                IconButton(
                  onPressed: () {},
                  icon: const Icon(Icons.more_vert, color: AppColors.accent),
                ),
              ],
            ),
          ),

          // Content
          if (reportText.isNotEmpty)
            Padding(
              padding: const EdgeInsets.symmetric(horizontal: 16),
              child: Text(reportText, style: AppTextStyles.body1),
            ),

          // Image
          if (imageUrl != null && imageUrl.isNotEmpty)
            Container(
              margin: const EdgeInsets.all(16),
              height: 200,
              width: double.infinity,
              decoration: BoxDecoration(
                color: AppColors.border.withValues(alpha: 0.3),
                borderRadius: BorderRadius.circular(12),
              ),
              child: ClipRRect(
                borderRadius: BorderRadius.circular(12),
                child: Image.network(
                  '$AppConstants.imageBaseUrl/$imageUrl',
                  fit: BoxFit.cover,
                  errorBuilder: (context, error, stackTrace) {
                    return const Center(
                      child: Icon(
                        Icons.image,
                        size: 50,
                        color: AppColors.accent,
                      ),
                    );
                  },
                  loadingBuilder: (context, child, loadingProgress) {
                    if (loadingProgress == null) return child;
                    return const Center(child: CircularProgressIndicator());
                  },
                ),
              ),
            ),

          // Actions
          Padding(
            padding: const EdgeInsets.all(16),
            child: Row(
              children: [
                _ActionButton(
                  icon: isLiked ? Icons.thumb_up : Icons.thumb_up_outlined,
                  label: '$likesCount',
                  onTap: onLike,
                  isActive: isLiked,
                ),
                const SizedBox(width: 24),
                _ActionButton(
                  icon: Icons.share_outlined,
                  label: '$sharesCount',
                  onTap: onShare,
                ),
                const Spacer(),
              ],
            ),
          ),
        ],
      ),
    );
  }
}

class _ActionButton extends StatelessWidget {
  final IconData icon;
  final String label;
  final VoidCallback onTap;
  final bool isActive;

  const _ActionButton({
    required this.icon,
    required this.label,
    required this.onTap,
    this.isActive = false,
  });

  @override
  Widget build(BuildContext context) {
    return GestureDetector(
      onTap: onTap,
      child: Row(
        children: [
          Icon(
            icon,
            size: 20,
            color: isActive ? AppColors.primary : AppColors.accent,
          ),
          const SizedBox(width: 6),
          Text(
            label,
            style: TextStyle(
              color: isActive ? AppColors.primary : AppColors.accent,
              fontWeight: isActive ? FontWeight.w600 : FontWeight.normal,
            ),
          ),
        ],
      ),
    );
  }
}
