import 'package:app/proto/chat.pb.dart';
import 'package:app/proto/chat.pbgrpc.dart';
import 'package:grpc/grpc.dart';
import 'package:intl/intl.dart';

class ChatService {
  ClientChannel channel;
  ChatClient stub;
  Message msg;
  List<ServerMessage> listMsg;

  ChatService() {
    msg = Message();
    listMsg = List<ServerMessage>();

    channel = ClientChannel('10.0.1.172',
        port: 9001,
        options:
            const ChannelOptions(credentials: ChannelCredentials.insecure()));

    stub = ChatClient(channel,
        options: CallOptions(timeout: Duration(seconds: 3600)));
  }

  Stream<List<ServerMessage>> updateMsg() async* {

    await for (final note in stub.sendMessage(streamIt())) {
      print(
          'Message received from grpc: ${note.message.username} ${note.timestamp} ${note.message.content}');
      if (note.message.idPerson != 1) {
        listMsg.add(note);
      }
      yield listMsg;
    }
  }

  Stream<Message> streamIt() async* {
    yield msg;
  }

  void add(String username, String content) {

    // update the class attribute msg
    msg.username = username;
    msg.content = content;
    msg.idPerson = 1;

    // this is a workaround that I found to add a new message to the list without "overwriting"**
    // all the other messages (https://github.com/git-hyagi/janaina/issues/15)
    // ** actually, it was not overwriting, the thing is, when I did srvMsg.message = msg
    // ** the srvMsg.message was receiving a reference of msg, not a copy of it.
    Message temp = new Message();
    temp.idPerson = 1;
    temp.username = username;
    temp.content  = content;

    var srvMsg = ServerMessage();
    srvMsg.timestamp = DateFormat('HH:mm').format(DateTime.now()).toString();
    srvMsg.message = temp;
    listMsg.add(srvMsg);
  }
}
