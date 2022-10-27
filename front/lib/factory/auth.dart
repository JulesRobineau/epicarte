import 'package:flutter_secure_storage/flutter_secure_storage.dart';

class AuthResponse {
  final String accessToken;
  final String refreshToken;
  final DateTime expiresAt;

  AuthResponse({
    required this.accessToken,
    required this.refreshToken,
    required this.expiresAt,
  });

  factory AuthResponse.fromJson(Map<String, dynamic> json) {
    return AuthResponse(
      accessToken: json['access_token'],
      refreshToken: json['refresh_token'],
      expiresAt: DateTime.parse(json['expires_in']),
    );
  }

  void persist(FlutterSecureStorage storage) {
    storage.write(key: "token", value: accessToken);
    storage.write(key: "refresh_token", value: refreshToken);
    storage.write(key: "expires_at", value: expiresAt.toString());
  }

}

class Account {
  final int id;
  final String email;
  final String firstName;
  final String lastName;
  final String role;
  final String username;

  Account({
    required this.id,
    required this.email,
    required this.firstName,
    required this.lastName,
    required this.role,
    required this.username,
  });

  factory Account.fromJson(Map<String, dynamic> json) {
    return Account(
      id: json['id'],
      email: json['email'],
      firstName: json['first_name'],
      lastName: json['last_name'],
      role: json['role'],
      username: json['username'],
    );
  }

  toJson() {
    return {
      'id': id,
      'email': email,
      'first_name': firstName,
      'last_name': lastName,
      'role': role,
      'username': username,
    };
  }
}
