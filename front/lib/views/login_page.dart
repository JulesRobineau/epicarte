import 'package:flutter/material.dart';
import 'package:flutter_secure_storage/flutter_secure_storage.dart';
import 'package:front/services/login.dart';
import 'package:front/views/home_page.dart';
import 'package:front/views/settings_page.dart';

class LoginPage extends StatefulWidget {
  const LoginPage({Key? key}) : super(key: key);
  static const String routeName = '/login';

  @override
  State<LoginPage> createState() => _LoginPageState();
}

class _LoginPageState extends State<LoginPage> {
  final TextEditingController _loginController = TextEditingController();
  final TextEditingController _passwordController = TextEditingController();
  final _storage = const FlutterSecureStorage();
  final _formKey = GlobalKey<FormState>();
  late final AuthService _authService;
  bool _visible = false;
  @override
  void initState() {
    _authService = AuthService();
    _authService.LoginFromStorage().then((value) =>
    {
      if (value != null) {
          Navigator.popAndPushNamed(context, HomePage.routeName)
      } else {
        setState(() {_isLoading = false;}),
      }
    });
    super.initState();
  }

  bool _isLoading = true;

  void _submit() {
    if (_formKey.currentState!.validate()) {
      _authService.Login(_loginController.text, _passwordController.text, true)
          .then((value) {
        if (value != null) {
          Navigator.popAndPushNamed(context, HomePage.routeName);
        }
      }).catchError((error, stackTrace) {
        ScaffoldMessenger.of(context).showSnackBar(SnackBar(
          content: Text(error.toString().replaceFirst("Exception:", "")),
        ));
      });
    }
  }

  @override
  Widget build(BuildContext context) {
    if (_isLoading) {
      return const Scaffold(body: Center(child: CircularProgressIndicator()));
    }

    return Scaffold(
      appBar: AppBar(
        backgroundColor: Colors.transparent,
        elevation: 0,
        actions: [
          IconButton(
              onPressed: () =>
              {
                Navigator.pushNamed(context, SettingsPage.routeName,
                    arguments: SettingsArguments(_storage))
              },
              icon: const Icon(Icons.settings, color: Colors.grey))
        ],
      ),
      body: Form(
        key: _formKey,
        child: Column(
          mainAxisAlignment: MainAxisAlignment.center,
          crossAxisAlignment: CrossAxisAlignment.center,
          children: [
            TextFormField(
              controller: _loginController,
              decoration: const InputDecoration(
                hintText: 'Username or Email',
                border: OutlineInputBorder(
                  borderRadius: BorderRadius.all(Radius.circular(10.0)),
                ),
                prefixIcon: Icon(Icons.person),
              ),
              validator: (value) {
                if (value == null || value.isEmpty) {
                  return 'Please enter some text';
                }

                if (value.length < 3) {
                  return 'Username must be at least 3 characters';
                }

                return null;
              },
            ),
            const SizedBox(height: 20),
            TextFormField(
              controller: _passwordController,
              obscureText: _visible,
              decoration: InputDecoration(
                hintText: 'Password',
                border: const OutlineInputBorder(
                  borderRadius: BorderRadius.all(Radius.circular(10.0)),
                ),
                prefixIcon: const Icon(Icons.lock),
                suffixIcon: IconButton(
                  icon:
                  Icon(_visible ? Icons.visibility_off : Icons.visibility),
                  onPressed: () =>
                      setState(() {
                        _visible = !_visible;
                      }),
                ),
              ),
              validator: (value) {
                if (value == null || value.isEmpty) {
                  return 'Please enter some text';
                }

                if (value.length < 8) {
                  return 'Password must be at least 8 characters';
                }

                return null;
              },
            ),
            const SizedBox(height: 20),
            ElevatedButton(
              onPressed: _submit,
              child: const Text('Login'),
            ),
          ],
        ),
      ),
    );
  }
}
