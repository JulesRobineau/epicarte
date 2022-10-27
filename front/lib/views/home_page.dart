import 'package:front/factory/class.dart';
import 'package:front/services/classe.dart';
import 'package:flutter/material.dart';
import 'package:front/services/login.dart';
import 'package:front/views/class_page.dart';
import 'package:front/views/login_page.dart';

class HomePage extends StatefulWidget {
  const HomePage({Key? key}) : super(key: key);
  static const String routeName = '/home';

  @override
  State<HomePage> createState() => _HomePageState();
}

class _HomePageState extends State<HomePage> {
  final ClassService _classeService = ClassService();
  final AuthService _authService = AuthService();

  @override
  Widget build(BuildContext context) {
    return Scaffold(
        appBar: AppBar(
          title: const Text('Classes'),
          actions: [
            IconButton(
                onPressed: () {
                  _authService.Logout()
                      .then((value) => {
                            Navigator.popAndPushNamed(
                                context, LoginPage.routeName)
                          })
                      .catchError((error, stackTrace) {
                    ScaffoldMessenger.of(context).showSnackBar(SnackBar(
                      content:
                          Text(error.toString().replaceFirst("Exception:", "")),
                    ));
                  });
                },
                icon: const Icon(Icons.logout))
          ],
        ),
        body: FutureBuilder(
            future: _classeService.GetClasses(),
            builder: (context, async) {
              if (async.hasError) {
                return Center(child: Text('Error ${async.error.toString()}'));
              }
              if (!async.hasData) {
                return const Center(child: CircularProgressIndicator());
              }
              Classes classes = async.data as Classes;
              return SizedBox(
                width: MediaQuery.of(context).size.width,
                height: MediaQuery.of(context).size.height * 0.9,
                child: ListView.builder(
                  padding: const EdgeInsets.all(8),
                  itemCount: classes.classes.length,
                  itemBuilder: (context, index) {
                    return Card(
                      elevation: 3,
                      shape: RoundedRectangleBorder(
                        borderRadius: BorderRadius.circular(5.0),
                      ),
                      child: ListTile(
                        title: Text(classes.classes[index].name),
                        subtitle:
                            Text('School year ${classes.classes[index].year}'),
                        trailing: const Icon(Icons.arrow_forward_ios),
                        onTap: () {
                          Navigator.pushNamed(context, ClassPage.routeName,
                              arguments: ClassArguments(
                                  classes.classes[index].id,
                                  classes.classes[index].name));
                        },
                      ),
                    );
                  },
                ),
              );
            }));
  }
}
