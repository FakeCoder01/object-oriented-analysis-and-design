#include "httplib.h"
#include "memento/MementoPattern.hpp"
#include "no_pattern/NoPatternEditor.hpp"
#include <iostream>

// globals
NoPatternEditor no_pattern_editor;
MementoEditor memento_editor;
Caretaker caretaker(&memento_editor);

int main() {
  httplib::Server svr;
  svr.set_default_headers(
      {{"Access-Control-Allow-Origin", "*"},
       {"Access-Control-Allow-Methods", "POST, GET, OPTIONS"},
       {"Access-Control-Allow-Headers", "Content-Type"}});

  svr.Options(R"(.*)", [](const httplib::Request &, httplib::Response &res) {
    res.status = 200;
  });

  // no pattern api
  svr.Get("/api/nopattern/shapes", [](const httplib::Request &,
                                      httplib::Response &res) {
    res.set_content(no_pattern_editor.getShapes().dump(), "application/json");
  });
  svr.Post("/api/nopattern/add",
           [](const httplib::Request &req, httplib::Response &res) {
             auto j = json::parse(req.body);
             no_pattern_editor.addShape(Shape::fromJson(j));
             res.set_content("{\"status\":\"ok\"}", "application/json");
           });
  svr.Post("/api/nopattern/undo",
           [](const httplib::Request &, httplib::Response &res) {
             no_pattern_editor.undo();
             res.set_content("{\"status\":\"ok\"}", "application/json");
           });
  svr.Post("/api/nopattern/redo",
           [](const httplib::Request &, httplib::Response &res) {
             no_pattern_editor.redo();
             res.set_content("{\"status\":\"ok\"}", "application/json");
           });
  svr.Post("/api/nopattern/clear",
           [](const httplib::Request &, httplib::Response &res) {
             no_pattern_editor.clear();
             res.set_content("{\"status\":\"ok\"}", "application/json");
           });

  // memento api
  svr.Get("/api/memento/shapes", [](const httplib::Request &,
                                    httplib::Response &res) {
    res.set_content(memento_editor.getShapes().dump(), "application/json");
  });
  svr.Post("/api/memento/add",
           [](const httplib::Request &req, httplib::Response &res) {
             caretaker.backup();
             auto j = json::parse(req.body);
             memento_editor.addShape(Shape::fromJson(j));
             res.set_content("{\"status\":\"ok\"}", "application/json");
           });
  svr.Post("/api/memento/undo",
           [](const httplib::Request &, httplib::Response &res) {
             caretaker.undo();
             res.set_content("{\"status\":\"ok\"}", "application/json");
           });
  svr.Post("/api/memento/redo",
           [](const httplib::Request &, httplib::Response &res) {
             caretaker.redo();
             res.set_content("{\"status\":\"ok\"}", "application/json");
           });
  svr.Post("/api/memento/clear",
           [](const httplib::Request &, httplib::Response &res) {
             caretaker.backup();
             memento_editor.clear();
             res.set_content("{\"status\":\"ok\"}", "application/json");
           });

  std::cout << "Starting server on port 8080..." << std::endl;
  svr.listen("0.0.0.0", 8080);
  return 0;
}
