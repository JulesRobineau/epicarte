class Student {
  int id;
  final String firstName;
  final String lastName;
  final String email;
  String role;

  Student(
      {this.id = 0,
      required this.firstName,
      required this.lastName,
      required this.email,
      this.role = "student"});

  @override
  String toString() {
    return "$firstName $lastName";
  }

  @override
  bool operator ==(Object other) {
    if (other.runtimeType != runtimeType) {
      return false;
    }
    return other is Student &&
        other.firstName == firstName &&
        other.lastName == lastName;
  }

  @override
  int get hashCode => Object.hash(email, firstName, lastName);

  factory Student.fromJson(Map<String, dynamic> json) {
    return Student(
      id: json['id'],
      firstName: json['first_name'],
      lastName: json['last_name'],
      email: json['email'],
      role: json['role'] ?? "student",
    );
  }
}

class Students {
  final students = <Student>[];

  Students.fromJson(List<dynamic> json) {
    for (final studentJson in json) {
      students.add(Student.fromJson(studentJson));
    }
  }
}
