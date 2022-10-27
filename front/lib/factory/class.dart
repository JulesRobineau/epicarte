class Class {
  final int id;
  final String name;
  final String year;

  Class({required this.id, required this.name, required this.year});

  factory Class.fromJson(Map<String, dynamic> json) {
    return Class(
      id: json['id'],
      name: json['name'],
      year: json['year'],
    );
  }
}

class Classes {
  final classes = <Class>[];

  Classes.fromJson(List<dynamic> json) {
    for (final classJson in json) {
      classes.add(Class.fromJson(classJson));
    }
  }
}
