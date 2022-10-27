import 'dart:convert';

import 'package:flutter/material.dart';
import 'package:front/factory/student.dart';
import 'package:web_socket_channel/web_socket_channel.dart';

class AddStudentFloatingButton extends StatelessWidget {
  const AddStudentFloatingButton(
      {super.key, required this.students, required this.channel});

  final List<Student>? students;
  final WebSocketChannel channel;

  @override
  Widget build(BuildContext context) {
    TextEditingController textFieldController = TextEditingController();

    return FloatingActionButton(
      onPressed: () => showDialog(
          context: context,
          builder: (BuildContext context) => AlertDialog(
                title: const Text('Ajouter un élève'),
                content: const Text('Veuillez entrer le nom de l\'élève'),
                actions: [
                  AutocompleteField(
                    students: students ?? [],
                    controller: textFieldController,
                  ),
                  Row(
                    mainAxisAlignment: MainAxisAlignment.end,
                    children: [
                      TextButton(
                        onPressed: () {
                          textFieldController.clear();
                          Navigator.pop(context);
                        },
                        child: const Text('Annuler'),
                      ),
                      TextButton(
                        onPressed: () {
                          channel.sink.add(json.encode({
                            'origin': 'mobile',
                            'action': 'addStudent',
                            'student': textFieldController.text
                          }));
                          textFieldController.clear();
                          Navigator.pop(context);
                        },
                        child: const Text('Ajouter'),
                      ),
                    ],
                  ),
                ],
              )),
      child: const Icon(Icons.add),
    );
  }
}

class AutocompleteField extends StatelessWidget {
  const AutocompleteField(
      {super.key, required this.students, required this.controller});

  final List<Student> students;
  final TextEditingController controller;

  static String _displayStringForOption(Student option) => option.toString();

  @override
  Widget build(BuildContext context) {
    return Autocomplete<Student>(
      initialValue: const TextEditingValue(text: ''),
      displayStringForOption: _displayStringForOption,
      fieldViewBuilder: (
        BuildContext context,
        TextEditingController textEditingController,
        FocusNode focusNode,
        VoidCallback onFieldSubmitted,
      ) {
        textEditingController.addListener(() {
          controller.text = textEditingController.text;
        });
        return TextField(
            controller: textEditingController,
            focusNode: focusNode,
            onSubmitted: (String value) {
              onFieldSubmitted();
            },
            decoration: const InputDecoration(
              hintText: 'Nom de l\'élève',
              border: OutlineInputBorder(
                borderRadius: BorderRadius.all(Radius.circular(10.0)),
              ),
              prefixIcon: Icon(Icons.person),
            ));
      },
      optionsBuilder: (TextEditingValue textEditingValue) {
        return students.where((Student option) {
          return option
              .toString()
              .toLowerCase()
              .contains(textEditingValue.text.toLowerCase());
        });
      },
      onSelected: (Student selection) {
        debugPrint('You just selected $selection');
      },
    );
  }
}
