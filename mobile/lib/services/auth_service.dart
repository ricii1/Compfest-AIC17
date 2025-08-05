import 'dart:convert';
import 'dart:io';
import 'package:http/http.dart' as http;
import '../utils/constants.dart';
import '../models/api_response.dart';

class AuthService {
  static const String baseUrl = '${AppConstants.apiBaseUrl}/user';

  // Singleton pattern for AuthService
  static final AuthService _instance = AuthService._internal();
  factory AuthService() => _instance;
  AuthService._internal();

  Future<ApiResponse<Map<String, dynamic>>> register({
    required String email,
    required String username,
    required String password,
  }) async {
    try {
      final url = Uri.parse(baseUrl);
      final body = {
        'email': email.trim(),
        'name': username.trim(),
        'password': password,
      };

      print('Sending register request to: $url');
      print('Request body: ${json.encode(body)}');

      final response = await http
          .post(
            url,
            headers: {
              'Content-Type': 'application/json',
              'Accept': 'application/json',
            },
            body: json.encode(body),
          )
          .timeout(const Duration(seconds: 10));

      print('Response status: ${response.statusCode}');
      print('Response body: ${response.body}');

      final responseData = json.decode(response.body) as Map<String, dynamic>;

      if (response.statusCode == 200 || response.statusCode == 201) {
        return ApiResponse.success(responseData);
      } else {
        final errorMessage =
            responseData['message'] ??
            responseData['error'] ??
            responseData['errors']?.toString() ??
            'Registration failed';
        return ApiResponse.error(errorMessage);
      }
    } on SocketException {
      return ApiResponse.error(
        'Unable to connect to server. Please check your internet connection and server availability.',
      );
    } on FormatException {
      return ApiResponse.error('Invalid response format from server.');
    } on http.ClientException {
      return ApiResponse.error('Network error occurred. Please try again.');
    } catch (e) {
      print('Registration error: $e');
      if (e.toString().contains('TimeoutException')) {
        return ApiResponse.error(
          'Connection timeout. Please check your internet connection.',
        );
      } else {
        return ApiResponse.error('Registration failed: ${e.toString()}');
      }
    }
  }

  Future<ApiResponse<Map<String, dynamic>>> login({
    required String username,
    required String password,
  }) async {
    try {
      final url = Uri.parse('$baseUrl/login');
      final body = {'username': username, 'password': password};

      print('Sending login request to: $url');
      print('Request body: ${json.encode(body)}');

      final response = await http
          .post(
            url,
            headers: {
              'Content-Type': 'application/json',
              'Accept': 'application/json',
            },
            body: json.encode(body),
          )
          .timeout(const Duration(seconds: 30));

      print('Response status: ${response.statusCode}');
      print('Response body: ${response.body}');

      final responseData = json.decode(response.body) as Map<String, dynamic>;

      if (response.statusCode == 200) {
        return ApiResponse.success(responseData);
      } else {
        final errorMessage =
            responseData['message'] ??
            responseData['error'] ??
            responseData['errors']?.toString() ??
            'Login failed';
        return ApiResponse.error(errorMessage);
      }
    } on SocketException {
      return ApiResponse.error(
        'Unable to connect to server. Please check your internet connection and server availability.',
      );
    } on FormatException {
      return ApiResponse.error('Invalid response format from server.');
    } on http.ClientException {
      return ApiResponse.error('Network error occurred. Please try again.');
    } catch (e) {
      print('Login error: $e');
      if (e.toString().contains('TimeoutException')) {
        return ApiResponse.error(
          'Connection timeout. Please check your internet connection.',
        );
      } else {
        return ApiResponse.error('Login failed: ${e.toString()}');
      }
    }
  }

  // Method untuk logout (jika diperlukan nanti)
  Future<void> logout() async {
    // Implement logout logic here
    // Clear stored tokens, user data, etc.
  }
}
