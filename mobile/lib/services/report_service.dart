import 'dart:convert';
import 'dart:io';
import '../models/report.dart';
import '../models/api_response.dart';
import 'api_service.dart';

class ReportService {
  static Future<ApiResponse<List<Report>>> getReports(String token) async {
    return await ApiService.get<List<Report>>(
      '/reports',
      token: token,
      fromJson: (json) =>
          (json['data'] as List).map((item) => Report.fromJson(item)).toList(),
    );
  }

  static Future<ApiResponse<Report>> createReport(
    String token,
    String text,
    File? image,
  ) async {
    if (image != null) {
      return await ApiService.postMultipart<Report>(
        '/reports',
        {'text': text},
        file: image,
        token: token,
        fromJson: (json) => Report.fromJson(json),
      );
    } else {
      return await ApiService.post<Report>(
        '/reports',
        {'text': text},
        token: token,
        fromJson: (json) => Report.fromJson(json),
      );
    }
  }
}
