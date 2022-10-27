import 'package:flutter/material.dart';
import 'package:flutter_secure_storage/flutter_secure_storage.dart';
import 'package:shared_preferences/shared_preferences.dart';

class SettingsArguments {
  late final FlutterSecureStorage? storage;

  SettingsArguments(this.storage);
}

class SettingsPage extends StatefulWidget {
  const SettingsPage({Key? key}) : super(key: key);
  static const String routeName = '/settings';

  @override
  State<SettingsPage> createState() => _SettingsPageState();
}

class _SettingsPageState extends State<SettingsPage> {
  final _formKey = GlobalKey<FormState>();
  final TextEditingController _serverController =
      TextEditingController(text: "http://localhost:8080");
  late final SharedPreferences prefs;

  @override
  void initState() {
    SharedPreferences.getInstance().then((value) => {
          setState(() {
            prefs = value;
            debugPrint(prefs.getString("server_url") ?? "null");
            _serverController.text =
                prefs.getString("server_url") ?? "http://localhost:8080";
          })
        });
    super.initState();
  }

  @override
  Widget build(BuildContext context) {
    final SettingsArguments args =
        ModalRoute.of(context)!.settings.arguments as SettingsArguments;
    args.storage ??= const FlutterSecureStorage();
    return Scaffold(
      appBar: AppBar(
        title: const Text('Settings'),
        leading: IconButton(
          icon: const Icon(Icons.arrow_back),
          onPressed: () => Navigator.of(context).pop(),
        ),
      ),
      body: Form(
        key: _formKey,
        child: Column(
          children: [
            const SizedBox(height: 15),
            TextFormField(
              controller: _serverController,
              decoration: const InputDecoration(
                labelText: 'Server Url',
                hintText: "Server URL",
                border: OutlineInputBorder(
                  borderRadius: BorderRadius.all(Radius.circular(10.0)),
                ),
              ),
              validator: (value) {
                if (value == null || value.isEmpty) {
                  return 'Server URL cannot be empty';
                }

                if (!value.startsWith("http://") &&
                    !value.startsWith("https://")) {
                  return 'Server URL must start with http:// or https://';
                }

                return null;
              },
            ),
            const SizedBox(height: 15),
            ElevatedButton(
              onPressed: () async {
                if (_formKey.currentState!.validate()) {
                  debugPrint("Saving server url: ${_serverController.text}");
                  await prefs.setString("server_url", _serverController.text);
                }
              },
              child: const Text('Update'),
            ),
            const SizedBox(height: 25),
            ElevatedButton(
                onPressed: () {
                  args.storage!.delete(key: "token");
                },
                style: ElevatedButton.styleFrom(backgroundColor: Colors.red),
                child: const Text('Supprime le token')),
            ElevatedButton(
              onPressed: () async {
                await prefs.clear();
              },
              style: ElevatedButton.styleFrom(backgroundColor: Colors.red),
              child: const Text('Suprrime les préférences'),
            ),
          ],
        ),
      ),
    );
  }
}
