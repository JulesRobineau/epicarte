import 'package:flutter/material.dart';
import 'package:front/views/class_page.dart';
import 'package:front/views/session_page.dart';
import 'package:front/views/home_page.dart';
import 'package:front/views/login_page.dart';
import 'package:front/views/settings_page.dart';
import 'package:front/views/student_page.dart';

void main() {
  runApp(const MyApp());
}

class MyApp extends StatelessWidget {
  const MyApp({super.key});

  // This widget is the root of your application.
  @override
  Widget build(BuildContext context) {
    return MaterialApp(
      title: 'Flutter Demo',
      theme: ThemeData(
        primarySwatch: Colors.blue,
      ),
      initialRoute: '/login',
      routes: {
        LoginPage.routeName: (context) => const LoginPage(),
        HomePage.routeName: (context) => const HomePage(),
        ClassPage.routeName: (context) => const ClassPage(),
        SessionPage.routeName: (context) => const SessionPage(),
        SettingsPage.routeName: (context) => const SettingsPage(),
        StudentsPage.routeName: (context) => const StudentsPage(),
      },
    );
  }
}
