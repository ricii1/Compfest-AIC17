class User {
  final String id;
  final String name;
  final String email;
  final String? phoneNumber;
  final String role;
  final String? imageUrl;
  final bool isVerified;

  User({
    required this.id,
    required this.name,
    required this.email,
    this.role = 'user',
    this.isVerified = false,
    this.phoneNumber,
    this.imageUrl,
  });
  factory User.fromJson(Map<String, dynamic> json) {
    return User(
      id: json['id'],
      email: json['email'],
      name: json['name'],
      imageUrl: json[''],
      role: json['role'] ?? 'user',
      isVerified: json['is_verified'] ?? false,
      phoneNumber: json['phone_number'],
    );
  }

  Map<String, dynamic> toJson() {
    return {
      'id': id,
      'email': email,
      'name': name,
      'avatar': imageUrl,
      'role': role,
      'is_verified': isVerified,
      'phone_number': phoneNumber,
    };
  }
}
