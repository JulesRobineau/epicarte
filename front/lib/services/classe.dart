import 'dart:convert';

import 'package:flutter_secure_storage/flutter_secure_storage.dart';
import 'package:front/factory/session.dart';
import 'package:front/factory/student.dart';
import 'package:shared_preferences/shared_preferences.dart';
import 'package:http/http.dart' as http;
import 'package:front/factory/class.dart';

class ClassService {
  final _storage = const FlutterSecureStorage();
  String _baseUrl = "http://localhost:8080/api/v1";

  Future<void> _setServerUrl() async {
    SharedPreferences prefs = await SharedPreferences.getInstance();
    String? su = prefs.getString("server_url");
    if (su != null) {
      _baseUrl = "$su/api/v1";
    }
  }

  String _bodyErrorFormat(Map error) {
    String message = "";
    for (var key in error.keys) {
      message += "$key: ${error[key]}. ";
    }

    return message;
  }

  Future<Sessions> GetClassSessions(int classId) async {
    await _setServerUrl();
    final token = await _storage.read(key: "token");
    if (token == null) {
      throw Exception("No token found");
    }
    final response = await http.get(
      Uri.parse('$_baseUrl/classes/$classId/sessions'),
      headers: <String, String>{
        'Content-Type': 'application/json; charset=UTF-8',
        'Authorization': 'Bearer $token',
      },
    );
    if (response.statusCode == 200) {
      return Sessions.fromJson(jsonDecode(response.body)["sessions"]);
    }
    throw Exception("Failed to get class sessions");
  }

  Future<Classes> GetClasses() async {
    await _setServerUrl();
    final token = await _storage.read(key: "token");
    if (token != null && token != "token") {
      final response = await http.get(
        Uri.parse('$_baseUrl/classes'),
        headers: <String, String>{
          'Content-Type': 'application/json; charset=UTF-8',
          'Authorization': 'Bearer $token',
        },
      );
      Map<String, dynamic> body = jsonDecode(response.body);
      if (response.statusCode == 200) {
        return Classes.fromJson(body["classes"] as List);
      }
      throw Exception(body["message"]);
    }
    throw Exception("Not authorized");
  }

  Future<Students> GetClassStudents(int classId) async {
    await _setServerUrl();
    final token = await _storage.read(key: "token");
    if (token != null && token != "token") {
      final response = await http.get(
        Uri.parse('$_baseUrl/classes/$classId'),
        headers: <String, String>{
          'Content-Type': 'application/json; charset=UTF-8',
          'Authorization': 'Bearer $token',
        },
      );
      Map<String, dynamic> body = jsonDecode(response.body);
      if (response.statusCode == 200) {
        return Students.fromJson(body["students"] as List);
      }
      throw Exception(body["message"]);
    }
    throw Exception("Not authorized");
  }

  Future<void> AddStudentToClass(
      int classId, Student student) async {
    await _setServerUrl();
    final token = await _storage.read(key: "token");
    if (token != null && token != "token") {
      final response = await http.post(
        Uri.parse('$_baseUrl/classes/$classId/students'),
        headers: <String, String>{
          'Content-Type': 'application/json; charset=UTF-8',
          'Authorization': 'Bearer $token',
        },
        body: jsonEncode({
          "student": {
            "first_name": student.firstName,
            "last_name": student.lastName,
            "email": student.email,
          }
        }),
      );
      var body = jsonDecode(response.body);
      if (response.statusCode == 201) {
        return;
      }
      if (response.statusCode == 400) {
        throw Exception(_bodyErrorFormat(body["error"]));
      }
      throw Exception(body["message"]);
    }
    throw Exception("Not authorized");
  }
}
