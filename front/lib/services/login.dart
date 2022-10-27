import 'dart:convert';

import 'package:flutter_secure_storage/flutter_secure_storage.dart';
import 'package:http/http.dart' as http;
import 'package:front/factory/auth.dart';
import 'package:shared_preferences/shared_preferences.dart';

class AuthService {
  late String _baseUrl;
  final FlutterSecureStorage _storage = const FlutterSecureStorage();

  AuthService() {
    _baseUrl = "http://localhost:8080/api/v1";
  }

  Future<void> _setServerUrl() async {
    SharedPreferences prefs = await SharedPreferences.getInstance();
    String? su = prefs.getString("server_url");
    if (su != null) {
      _baseUrl = "$su/api/v1";
    }
  }

  Future<Account?> LoginFromStorage() async {
    await _setServerUrl();
    final token = await _storage.read(key: "token");
    if (token != null && token != "token") {
      return GetAccount(token);
    }
    return null;
  }

  Future<AuthResponse?> Login(String login, String password,
      [bool? withPersistence]) async {
    await _setServerUrl();
    final response = await http.post(
      Uri.parse('$_baseUrl/auth/login'),
      headers: <String, String>{
        'Content-Type': 'application/json; charset=UTF-8',
      },
      body: jsonEncode(<String, String>{
        'username': login,
        'password': password,
      }),
    );
    Map<String, dynamic> body = jsonDecode(response.body);
    if (response.statusCode == 202) {
      final AuthResponse authResponse =
          AuthResponse.fromJson(body);
      if (withPersistence ?? false) {
        authResponse.persist(_storage);
      }

      return authResponse;
    }
    if ((body["message"] as String) == "record not found") {
      throw Exception("Login or username is incorrect");
    }

    throw Exception(body["message"]);
  }

  Future<Account?> GetAccount(String accessToken) async {
    await _setServerUrl();
    final response = await http.get(
      Uri.parse('$_baseUrl/auth/account'),
      headers: <String, String>{
        'Content-Type': 'application/json; charset=UTF-8',
        'Authorization': 'Bearer $accessToken',
      },
    );

    if (response.statusCode == 200) {
      return Account.fromJson(jsonDecode(response.body));
    }

    return null;
  }

  Future<AuthResponse?> RefreshToken(String refreshToken) async {
    await _setServerUrl();
    final response = await http.post(
      Uri.parse('$_baseUrl/auth/refresh-token'),
      headers: <String, String>{
        'Content-Type': 'application/json; charset=UTF-8',
      },
      body: jsonEncode(<String, String>{
        'refresh_token': refreshToken,
      }),
    );

    if (response.statusCode == 202) {
      return AuthResponse.fromJson(jsonDecode(response.body));
    }

    return null;
  }

  Future<void> Logout() async {
    await _setServerUrl();
    final refreshToken = await _storage.read(key: "refresh_token");
    final response = await http.delete(
      Uri.parse('$_baseUrl/auth/logout'),
      headers: <String, String>{
        'Content-Type': 'application/json; charset=UTF-8',
        'Authorization': 'Bearer $refreshToken',
      },
    );

    if (response.statusCode == 204) {
      await _storage.delete(key: "token");
      await _storage.delete(key: "refresh_token");
      await _storage.delete(key: "expires_at");
      return;
    }

    throw Exception("Logout failed");
  }
}
