import 'dart:io';

class Report {
  final String? id;
  final String text;
  final String? imageUrl;
  final File? imageFile;
  final String userId;
  final String userName;
  final DateTime createdAt;

  Report({
    this.id,
    required this.text,
    this.imageUrl,
    this.imageFile,
    required this.userId,
    required this.userName,
    required this.createdAt,
  });

  factory Report.fromJson(Map<String, dynamic> json) {
    return Report(
      id: json['id'],
      text: json['text'],
      imageUrl: json['image_url'],
      userId: json['user_id'],
      userName: json['user_name'],
      createdAt: DateTime.parse(json['created_at']),
    );
  }
  Map<String, dynamic> toJson() {
    return {
      'text': text,
      'image': imageFile != null ? 'base64_encoded_image' : null,
    };
  }
}
