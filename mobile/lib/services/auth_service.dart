import '../models/user.dart';
import '../models/api_response.dart';
import 'api_service.dart';

class AuthService {
  static Future<ApiResponse<String>> login(
    String email,
    String password,
  ) async {
    final response = await ApiService.post<Map<String, dynamic>>(
      '/auth/login',
      {'email': email, 'password': password},
    );

    if (response.isSuccess) {
      return ApiResponse.success(response.data!['token']);
    } else {
      return ApiResponse.error(response.message);
    }
  }

  static Future<ApiResponse<User>> getProfile(String token) async {
    return await ApiService.get<User>(
      '/auth/profile',
      token: token,
      fromJson: (json) => User.fromJson(json),
    );
  }
}
