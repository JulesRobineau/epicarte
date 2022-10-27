import 'package:flutter_secure_storage/flutter_secure_storage.dart';
import 'package:shared_preferences/shared_preferences.dart';
import 'package:web_socket_channel/web_socket_channel.dart';

class SessionWsService {
  final FlutterSecureStorage _storage = const FlutterSecureStorage();
  String _baseUrl = 'ws://localhost:8080';

  Future<void> _setServerUrl() async {
    SharedPreferences prefs = await SharedPreferences.getInstance();
    String? su = prefs.getString("server_url");
    if (su != null) {
      su = su.replaceFirst("http", "ws");
      _baseUrl = "$su/api/v1";
    }
  }

  Future<WebSocketChannel> GetSesssionWs(String id, String password) async {
    await _setServerUrl();
    return WebSocketChannel.connect(
      Uri.parse('$_baseUrl/ws/sessions/$id/join?password=$password'),
    );
  }
}
