import 'package:flutter/material.dart';
// import 'screens/auth/login_screen.dart';
import 'utils/constants.dart';
import 'widgets/auth_guard.dart';
import 'screens/home/home_screen.dart';

void main() {
  runApp(const ReportApp());
}

class ReportApp extends StatelessWidget {
  const ReportApp({super.key});

  @override
  Widget build(BuildContext context) {
    return MaterialApp(
      title: AppStrings.appName,
      theme: ThemeData(
        colorScheme: ColorScheme.fromSeed(seedColor: AppColors.primary),
        useMaterial3: true,
        scaffoldBackgroundColor: AppColors.background,
        appBarTheme: const AppBarTheme(
          backgroundColor: AppColors.primary,
          foregroundColor: Colors.white,
          elevation: 0,
        ),
      ),
      home: const AuthGuard(child: HomeScreen(token: '')),
      debugShowCheckedModeBanner: false,
    );
  }
}
