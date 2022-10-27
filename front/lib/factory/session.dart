class Session {
  final String id;
  final bool isClosed;
  final DateTime createdAt;
  final DateTime updatedAt;



  Session({
    required this.id,
    required this.isClosed,
    required this.createdAt,
    required this.updatedAt,
  });

  factory Session.fromJson(Map<String, dynamic> json) {
    return Session(
      id: json['id'],
      isClosed: json['is_closed'],
      createdAt: DateTime.parse(json['created_at']),
      updatedAt: DateTime.parse(json['updated_at']),
    );
  }
}

class Sessions {
  final sessions = <Session>[];

  Sessions.fromJson(List<dynamic> json) {
    for (final sessionJson in json) {
      sessions.add(Session.fromJson(sessionJson));
    }
  }
}