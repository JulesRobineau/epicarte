import 'package:flutter/material.dart';
import 'package:front/factory/student.dart';
import 'package:front/services/classe.dart';

class StudentsArguments {
  final int classId;

  StudentsArguments(this.classId);
}

class StudentsPage extends StatefulWidget {
  const StudentsPage({Key? key}) : super(key: key);
  static const String routeName = '/students';

  @override
  State<StudentsPage> createState() => _StudentListPageState();
}

class _StudentListPageState extends State<StudentsPage> {
  final ClassService _classService = ClassService();

  @override
  Widget build(BuildContext context) {
    final StudentsArguments args =
        ModalRoute.of(context)!.settings.arguments as StudentsArguments;

    return Scaffold(
      appBar: AppBar(
        title: const Text('Students'),
      ),
      body: FutureBuilder(
          future: _classService.GetClassStudents(args.classId),
          builder: (context, async) {
            if (async.hasError) {
              return Center(child: Text('Error ${async.error.toString()}'));
            }
            if (!async.hasData) {
              return const Center(child: CircularProgressIndicator());
            }

            Students s = (async.data as Students);
            if (s.students.isEmpty) {
              return const Center(child: Text('No students'));
            }

            return SizedBox(
              width: MediaQuery.of(context).size.width,
              height: MediaQuery.of(context).size.height * 0.9,
              child: ListView.builder(
                padding: const EdgeInsets.all(8),
                itemCount: s.students.length,
                itemBuilder: (context, index) {
                  return Card(
                    child: ListTile(
                      title: Text(
                          "${s.students[index].firstName} ${s.students[index].lastName}"),
                      subtitle: Text(s.students[index].email),
                    ),
                  );
                },
              ),
            );
          }),
      floatingActionButton: FloatingActionButton(
        onPressed: () {
          showDialog(
              context: context,
              builder: (context) => _addStudentDialog(context, args.classId));
        },
        child: const Icon(Icons.person_add),
      ),
    );
  }

  Widget _addStudentDialog(BuildContext context, int classId) {
    TextEditingController firstNameController = TextEditingController();
    TextEditingController lastNameController = TextEditingController();
    TextEditingController emailController = TextEditingController();

    return AlertDialog(
      title: const Text('Add student'),
      content: Column(
        mainAxisSize: MainAxisSize.min,
        children: [
          TextField(
            controller: firstNameController,
            decoration: const InputDecoration(
              border: OutlineInputBorder(),
              labelText: 'First name',
            ),
          ),
          const SizedBox(height: 10),
          TextField(
            controller: lastNameController,
            decoration: const InputDecoration(
              border: OutlineInputBorder(),
              labelText: 'Last name',
            ),
          ),
          const SizedBox(height: 10),
          TextField(
            controller: emailController,
            decoration: const InputDecoration(
              border: OutlineInputBorder(),
              labelText: 'Email',
            ),
          ),
        ],
      ),
      actions: [
        TextButton(
          onPressed: () {
            Navigator.pop(context);
          },
          child: const Text('Cancel'),
        ),
        TextButton(
          onPressed: () async {
            await _classService.AddStudentToClass(
                    classId,
                    Student(
                        firstName: firstNameController.text,
                        lastName: lastNameController.text,
                        email: emailController.text))
                .then((value) {
                  Navigator.pop(context);
                  setState(() {});
                })
                .catchError((error, stackTrace) {
              debugPrint(error.toString());
              ScaffoldMessenger.of(context).showSnackBar(SnackBar(
                content: Text(error.toString().replaceFirst("Exception:", "")),
              ));
            });
          },
          child: const Text('Add'),
        ),
      ],
    );
  }
}
