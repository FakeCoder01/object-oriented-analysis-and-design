#include "MementoPattern.hpp"

Memento::Memento(const std::vector<Shape> &state) : state(state) {}

std::vector<Shape> Memento::getState() const { return state; }

void MementoEditor::addShape(const Shape &s) { shapes.push_back(s); }
void MementoEditor::clear() { shapes.clear(); }

std::unique_ptr<Memento> MementoEditor::save() {
  return std::make_unique<Memento>(shapes);
}

void MementoEditor::restore(const Memento *m) {
  if (m)
    shapes = m->getState();
}

json MementoEditor::getShapes() {
  json arr = json::array();
  for (auto &s : shapes)
    arr.push_back(s.toJson());
  return arr;
}

Caretaker::Caretaker(MementoEditor *ed) : editor(ed) {}

void Caretaker::backup() {
  undo_stack.push(editor->save());
  while (!redo_stack.empty())
    redo_stack.pop();
}

void Caretaker::undo() {
  if (!undo_stack.empty()) {
    redo_stack.push(editor->save());
    editor->restore(undo_stack.top().get());
    undo_stack.pop();
  }
}

void Caretaker::redo() {
  if (!redo_stack.empty()) {
    undo_stack.push(editor->save());
    editor->restore(redo_stack.top().get());
    redo_stack.pop();
  }
}
