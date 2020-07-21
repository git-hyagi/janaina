import 'package:app/chat_service.dart';
import 'package:app/proto/chat.pb.dart';
import 'package:flutter/cupertino.dart';
import 'package:flutter/material.dart';

void main() => runApp(MyApp());

class MyApp extends StatelessWidget {
  @override

  Widget build(BuildContext context) {
    return MaterialApp(
      title: 'Telemedicine',
      home: MyCustomForm(),
      debugShowCheckedModeBanner: false,
    );
  }
}

// Define a custom Form widget.
class MyCustomForm extends StatefulWidget {
  @override
  _MyCustomFormState createState() => _MyCustomFormState();
}

// Define a corresponding State class.
// This class holds the data related to the Form.
class _MyCustomFormState extends State<MyCustomForm> {
  // Create a text controller and use it to retrieve the current value
  // of the TextField.
  final myController = TextEditingController();
  ChatService streamGRPC = ChatService();
  ScrollController _scrollController = new ScrollController();

  @override
  void dispose() {
    // Clean up the controller when the widget is disposed.
    myController.dispose();
    super.dispose();
  }

  @override
  Widget build(BuildContext context) {

    return Scaffold(
      appBar: AppBar(
        title: Text('Telemedicine'),
      ),
      resizeToAvoidBottomPadding: false,
      resizeToAvoidBottomInset: true,
      body: Column(
        mainAxisAlignment: MainAxisAlignment.end,
        children: <Widget>[
          StreamBuilder<List<ServerMessage>>(
              stream:  streamGRPC.updateMsg(),
              builder: (context, snapshot){
                if (snapshot.hasError) {
                  return Text("Error: ${snapshot.error}");
                }
                if (snapshot.hasData) {
                  List<ServerMessage> mensagem = snapshot.data;
                  return Expanded (child: ListView.builder(
                    scrollDirection: Axis.vertical,
                    itemCount: mensagem.length,
                    controller: _scrollController,
                    shrinkWrap: true,
                    itemBuilder: (BuildContext context, int index){
                      return  ListTile(
                        title: Text('${mensagem[index].getField(2).username}: ${mensagem[index].getField(2).content}', style: TextStyle( color: Colors.blue[600])),
                        subtitle: Text(mensagem[index].getField(1).toString()),
                      );
                    } ,
                  ),
                  );}
                return Container();
              }
          ),
          Container(
            alignment: Alignment.bottomCenter,
            padding: const EdgeInsets.all(10.0),
            child: TextField(
              controller: myController,
              decoration: InputDecoration(border: OutlineInputBorder()),
            ),
          ),
        ],
      ),
      floatingActionButton: FloatingActionButton(
        onPressed: () {
          streamGRPC.add('android',myController.text);
          setState(() {});
          myController.clear();
          _scrollController.jumpTo(_scrollController.position.maxScrollExtent);
        },
        tooltip: 'Show me the value!',
        child: Icon(Icons.text_fields),
      ),
    );
  }
}