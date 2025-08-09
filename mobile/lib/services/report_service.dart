import 'dart:convert';
import 'dart:io';
import 'package:http/http.dart' as http;
import '../utils/constants.dart';
import '../models/api_response.dart';
import 'token_service.dart';

class ReportService {
  static const String baseUrl = '${AppConstants.apiBaseUrl}/reports';

  // Singleton pattern for ReportService
  static final ReportService _instance = ReportService._internal();
  factory ReportService() => _instance;
  ReportService._internal();

  // Get all reports with pagination
  Future<ApiResponse<Map<String, dynamic>>> getAllReports({
    int page = 1,
    int perPage = 10,
    String? search,
  }) async {
    try {
      final token = await TokenService.getToken();
      if (token == null) {
        return ApiResponse.error('User not authenticated');
      }

      final queryParams = {
        'page': page.toString(),
        'per_page': perPage.toString(),
        if (search != null && search.isNotEmpty) 'search': search,
      };

      final uri = Uri.parse(baseUrl).replace(queryParameters: queryParams);

      print('Fetching reports from: $uri');

      final response = await http
          .get(
            uri,
            headers: {
              'Content-Type': 'application/json',
              'Accept': 'application/json',
              'Authorization': 'Bearer $token',
            },
          )
          .timeout(const Duration(seconds: 30));

      print('Response status: ${response.statusCode}');
      print('Response body: ${response.body}');

      if (response.body.isEmpty) {
        return ApiResponse.error('Server returned empty response');
      }

      final responseData = json.decode(response.body) as Map<String, dynamic>;

      if (response.statusCode == 200) {
        // Return the full response data including pagination info
        return ApiResponse.success(responseData);
      } else {
        final errorMessage =
            responseData['message'] ??
            responseData['error'] ??
            'Failed to fetch reports';
        return ApiResponse.error(errorMessage);
      }
    } on SocketException {
      return ApiResponse.error(
        'Unable to connect to server. Please check your internet connection.',
      );
    } on FormatException catch (e) {
      return ApiResponse.error('Invalid response format: ${e.message}');
    } on http.ClientException catch (e) {
      return ApiResponse.error('Network error: ${e.message}');
    } catch (e) {
      print('Fetch reports error: $e');
      return ApiResponse.error('Failed to fetch reports: ${e.toString()}');
    }
  }

  // Get user's reports with pagination
  Future<ApiResponse<Map<String, dynamic>>> getUserReports({
    int page = 1,
    int perPage = 10,
    String? search,
  }) async {
    try {
      final token = await TokenService.getToken();
      if (token == null) {
        return ApiResponse.error('User not authenticated');
      }

      final userData = await TokenService.getUserData();
      final userId = userData['userId'];

      if (userId == null || userId.isEmpty) {
        return ApiResponse.error('User ID not found');
      }

      final queryParams = {
        'page': page.toString(),
        'per_page': perPage.toString(),
        if (search != null && search.isNotEmpty) 'search': search,
      };

      final uri = Uri.parse(
        '$baseUrl/user/$userId',
      ).replace(queryParameters: queryParams);

      print('Fetching user reports from: $uri');

      final response = await http
          .get(
            uri,
            headers: {
              'Content-Type': 'application/json',
              'Accept': 'application/json',
              'Authorization': 'Bearer $token',
            },
          )
          .timeout(const Duration(seconds: 30));

      print('Response status: ${response.statusCode}');
      print('Response body: ${response.body}');

      if (response.body.isEmpty) {
        return ApiResponse.error('Server returned empty response');
      }

      final responseData = json.decode(response.body) as Map<String, dynamic>;

      if (response.statusCode == 200) {
        return ApiResponse.success(responseData);
      } else {
        final errorMessage =
            responseData['message'] ??
            responseData['error'] ??
            'Failed to fetch user reports';
        return ApiResponse.error(errorMessage);
      }
    } on SocketException {
      return ApiResponse.error(
        'Unable to connect to server. Please check your internet connection.',
      );
    } on FormatException catch (e) {
      return ApiResponse.error('Invalid response format: ${e.message}');
    } on http.ClientException catch (e) {
      return ApiResponse.error('Network error: ${e.message}');
    } catch (e) {
      print('Fetch user reports error: $e');
      return ApiResponse.error('Failed to fetch user reports: ${e.toString()}');
    }
  }

  // Like/Unlike report
  Future<ApiResponse<Map<String, dynamic>>> toggleLike(String reportId) async {
    try {
      final token = await TokenService.getToken();
      if (token == null) {
        return ApiResponse.error('User not authenticated');
      }

      final url = Uri.parse('$baseUrl/$reportId/like');

      print('Toggling like for report: $url');

      final response = await http
          .post(
            url,
            headers: {
              'Content-Type': 'application/json',
              'Accept': 'application/json',
              'Authorization': 'Bearer $token',
            },
          )
          .timeout(const Duration(seconds: 30));

      print('Response status: ${response.statusCode}');
      print('Response body: ${response.body}');

      if (response.body.isEmpty) {
        return ApiResponse.error('Server returned empty response');
      }

      final responseData = json.decode(response.body) as Map<String, dynamic>;

      if (response.statusCode == 200 || response.statusCode == 201) {
        return ApiResponse.success(responseData);
      } else {
        final errorMessage =
            responseData['message'] ??
            responseData['error'] ??
            'Failed to toggle like';
        return ApiResponse.error(errorMessage);
      }
    } on SocketException {
      return ApiResponse.error(
        'Unable to connect to server. Please check your internet connection.',
      );
    } on FormatException catch (e) {
      return ApiResponse.error('Invalid response format: ${e.message}');
    } on http.ClientException catch (e) {
      return ApiResponse.error('Network error: ${e.message}');
    } catch (e) {
      print('Toggle like error: $e');
      return ApiResponse.error('Failed to toggle like: ${e.toString()}');
    }
  }

  // Create new report
  Future<ApiResponse<Map<String, dynamic>>> createReport({
    required String text,
    String? imageBase64,
    String? location,
  }) async {
    try {
      final token = await TokenService.getToken();
      if (token == null) {
        return ApiResponse.error('User not authenticated');
      }

      final url = Uri.parse(baseUrl);
      final body = {
        'text': text.trim(),
        if (imageBase64 != null) 'image': imageBase64,
        if (location != null) 'location': location,
      };

      print('Creating report at: $url');
      print('Request body keys: ${body.keys}');

      final response = await http
          .post(
            url,
            headers: {
              'Content-Type': 'application/json',
              'Accept': 'application/json',
              'Authorization': 'Bearer $token',
            },
            body: json.encode(body),
          )
          .timeout(const Duration(seconds: 30));

      print('Response status: ${response.statusCode}');
      print('Response body: ${response.body}');

      if (response.body.isEmpty) {
        return ApiResponse.error('Server returned empty response');
      }

      final responseData = json.decode(response.body) as Map<String, dynamic>;

      if (response.statusCode == 200 || response.statusCode == 201) {
        return ApiResponse.success(responseData);
      } else {
        final errorMessage =
            responseData['message'] ??
            responseData['error'] ??
            'Failed to create report';
        return ApiResponse.error(errorMessage);
      }
    } on SocketException {
      return ApiResponse.error(
        'Unable to connect to server. Please check your internet connection.',
      );
    } on FormatException catch (e) {
      return ApiResponse.error('Invalid response format: ${e.message}');
    } on http.ClientException catch (e) {
      return ApiResponse.error('Network error: ${e.message}');
    } catch (e) {
      print('Create report error: $e');
      return ApiResponse.error('Failed to create report: ${e.toString()}');
    }
  }
}
