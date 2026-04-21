#include "NoPatternEditor.hpp"

void NoPatternEditor::addShape(const Shape &s) {
  history.push(current_shapes);
  while (!redo_stack.empty())
    redo_stack.pop();
  current_shapes.push_back(s);
}

void NoPatternEditor::undo() {
  if (!history.empty()) {
    redo_stack.push(current_shapes);
    current_shapes = history.top();
    history.pop();
  }
}

void NoPatternEditor::redo() {
  if (!redo_stack.empty()) {
    history.push(current_shapes);
    current_shapes = redo_stack.top();
    redo_stack.pop();
  }
}

void NoPatternEditor::clear() {
  history.push(current_shapes);
  while (!redo_stack.empty())
    redo_stack.pop();
  current_shapes.clear();
}

json NoPatternEditor::getShapes() {
  json arr = json::array();
  for (auto &s : current_shapes)
    arr.push_back(s.toJson());
  return arr;
}
