import 'package:flutter/material.dart';

class AppConstants {
  static const String apiBaseUrl =
      'https://backend-aic-115083484414.asia-southeast2.run.app/api';
  static const String imageBaseUrl =
      'https://backend-aic-115083484414.asia-southeast2.run.app/assets';
  static const int maxTextLength = 280;
}

class AppColors {
  static const Color primary = Color(0xFF1DA1F2);
  static const Color secondary = Color(0xFF14171A);
  static const Color accent = Color(0xFF657786);
  static const Color background = Color(0xFFF7F9FA);
  static const Color cardBackground = Colors.white;
  static const Color error = Color(0xFFE0245E);
  static const Color success = Color(0xFF17BF63);
  static const Color warning = Color(0xFFFFAD1F);
  static const Color textPrimary = Color(0xFF14171A);
  static const Color textSecondary = Color(0xFF657786);
  static const Color border = Color(0xFFE1E8ED);
}

class AppStrings {
  static const String appName = 'Report App';
  static const String login = 'Login';
  static const String email = 'Email';
  static const String password = 'Password';
  static const String createReport = 'Create Report';
  static const String whatsHappening = "What's happening in your area?";
  static const String dontHaveAccount = "Don't have an account? Register";
  static const String welcomeBack = "Welcome Back";
  static const String loginSubtitle = "Sign in to continue reporting issues";
  static const String reports = 'Reports';
  static const String post = 'Post';
  static const String camera = 'Camera';
  static const String gallery = 'Gallery';
  static const String addImage = 'Add Image';
  static const String addLocation = 'Add Location';
  static const String like = 'Like';
  static const String comment = 'Comment';
  static const String share = 'Share';
  static const String failedToConnectToServer =
      'Failed to connect to the server. Please check your internet connection or try again later.';
  static const String networkError =
      'Network error occurred. Please try again.';
  static const String timeoutError =
      'Connection timeout. Please check your internet connection.';
  static const String invalidResponseFormat =
      'Invalid response format from server. Please contact support.';
}

class AppTextStyles {
  static const TextStyle headline1 = TextStyle(
    fontSize: 28,
    fontWeight: FontWeight.bold,
    color: AppColors.textPrimary,
  );

  static const TextStyle headline2 = TextStyle(
    fontSize: 24,
    fontWeight: FontWeight.bold,
    color: AppColors.textPrimary,
  );

  static const TextStyle subtitle1 = TextStyle(
    fontSize: 16,
    color: AppColors.textSecondary,
  );

  static const TextStyle body1 = TextStyle(
    fontSize: 16,
    color: AppColors.textPrimary,
  );

  static const TextStyle body2 = TextStyle(
    fontSize: 14,
    color: AppColors.textSecondary,
  );

  static const TextStyle caption = TextStyle(
    fontSize: 12,
    color: AppColors.textSecondary,
  );
}
