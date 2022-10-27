import 'package:flutter/material.dart';
import 'package:front/factory/session.dart';
import 'package:front/services/classe.dart';
import 'package:front/services/session.dart';
import 'package:front/views/session_page.dart';
import 'package:front/views/student_page.dart';
import 'package:web_socket_channel/web_socket_channel.dart';

class ClassArguments {
  final int classId;
  final String className;

  ClassArguments(this.classId, this.className);
}

class ClassPage extends StatefulWidget {
  const ClassPage({Key? key}) : super(key: key);
  static const String routeName = '/class';

  @override
  State<ClassPage> createState() => _ClassPageState();
}

class _ClassPageState extends State<ClassPage> {
  final ClassService _classService = ClassService();
  final SessionWsService _sessionWsService = SessionWsService();

  @override
  Widget build(BuildContext context) {
    final args = ModalRoute.of(context)!.settings.arguments as ClassArguments;
    final TextEditingController passwordController = TextEditingController();

    return Scaffold(
        appBar: AppBar(
          title: Text('Class ${args.className}'),
          actions: [
            IconButton(
                onPressed: () {
                  Navigator.pushNamed(context, StudentsPage.routeName,
                      arguments: StudentsArguments(args.classId));
                },
                icon: const Icon(Icons.person)),
          ],
        ),
        body: FutureBuilder(
            future: _classService.GetClassSessions(args.classId),
            builder: (context, async) {
              if (async.hasError) {
                return Center(child: Text('Error ${async.error.toString()}'));
              }
              if (!async.hasData) {
                return const Center(child: CircularProgressIndicator());
              }

              Sessions sessions = (async.data as Sessions);
              return SizedBox(
                width: MediaQuery.of(context).size.width,
                height: MediaQuery.of(context).size.height * 0.9,
                child: ListView.builder(
                  padding: const EdgeInsets.all(8),
                  itemCount: sessions.sessions.length,
                  itemBuilder: (context, index) {
                    return Card(
                      elevation: 3,
                      shape: RoundedRectangleBorder(
                        borderRadius: BorderRadius.circular(5.0),
                      ),
                      child: ListTile(
                        title: Text('Session: ${sessions.sessions[index].id}'),
                        trailing: const Icon(Icons.play_circle),
                        onTap: () => showDialog(
                          context: context,
                          builder: (BuildContext context) => AlertDialog(
                            title: Text('Join Session $index'),
                            content: const Text('Session\'s password ?'),
                            actions: [
                              TextField(
                                obscureText: true,
                                controller: passwordController,
                                decoration: const InputDecoration(
                                  hintText: 'Password',
                                  border: OutlineInputBorder(
                                    borderRadius:
                                        BorderRadius.all(Radius.circular(10.0)),
                                  ),
                                  prefixIcon: Icon(Icons.lock),
                                ),
                              ),
                              Row(
                                mainAxisAlignment: MainAxisAlignment.end,
                                children: [
                                  TextButton(
                                    onPressed: () {
                                      passwordController.clear();
                                      Navigator.of(context).pop();
                                    },
                                    child: const Text('Cancel'),
                                  ),
                                  TextButton(
                                    onPressed: () {
                                      _sessionWsService.GetSesssionWs(
                                              sessions.sessions[index].id,
                                              passwordController.text)
                                          .then((channel) => {
                                                _classService.GetClassStudents(
                                                        args.classId)
                                                    .then((students) =>
                                                        Navigator.pushNamed(
                                                            context,
                                                            SessionPage
                                                                .routeName,
                                                            arguments:
                                                                SessionArguments(
                                                                    channel,
                                                                    args.classId,
                                                                    students)))
                                              });
                                      passwordController.clear();
                                      Navigator.of(context).pop();
                                    },
                                    child: const Text('Join'),
                                  ),
                                ],
                              ),
                            ],
                          ),
                        ),
                      ),
                    );
                  },
                ),
              );
            }));
  }
}
