import 'package:flutter/material.dart';
import 'package:front/factory/student.dart';
import 'package:front/services/classe.dart';
import 'package:front/widgets/session/session_floating_button.dart';
import 'package:web_socket_channel/web_socket_channel.dart';

class SessionArguments {
  final WebSocketChannel channel;
  final int classId;
  final Students students;

  SessionArguments(this.channel, this.classId, this.students);
}

class SessionPage extends StatefulWidget {
  const SessionPage({Key? key}) : super(key: key);
  static const String routeName = '/session';

  @override
  State<SessionPage> createState() => _SessionPageState();
}

class _SessionPageState extends State<SessionPage> {
  late WebSocketChannel channel;

  @override
  void initState() {
    super.initState();
  }

  Widget _alertDialog(BuildContext context, String message) {
    return AlertDialog(
      title: const Text('Error'),
      content: Text(message),
      actions: [
        TextButton(
          onPressed: () {
            Navigator.pop(context);
          },
          child: const Text('Ok'),
        )
      ],
    );
  }

  @override
  Widget build(BuildContext context) {
    final args = ModalRoute.of(context)!.settings.arguments as SessionArguments;
    channel = args.channel;
    return Scaffold(
      appBar: AppBar(
        title: const Text('Sessions'),
        leading: IconButton(
          icon: const Icon(Icons.arrow_back),
          onPressed: () {
            Navigator.pop(context);
          },
        ),
      ),
      body: StreamBuilder(
          stream: channel.stream,
          builder: (context, async) {
            if (async.hasError) {
              debugPrint(async.error.toString());
              return _alertDialog(context, "Incorrect password");
            }

            if (async.hasData) {
              return Center(child: Text('Data: ${async.data}'));
            }

            return Center(
                child: Container(
                    width: 200,
                    height: 200,
                    decoration: BoxDecoration(
                      image: DecorationImage(
                        image: const AssetImage('images/nfc-logo.jpg'),
                        colorFilter: ColorFilter.mode(
                            Colors.black.withOpacity(0.2), BlendMode.dstATop),
                        fit: BoxFit.cover,
                      ),
                    )));
          }),
      floatingActionButton:
          AddStudentFloatingButton(students: args.students.students, channel: channel),
    );
  }

  @override
  void dispose() {
    channel.sink.close();
    super.dispose();
  }
}
