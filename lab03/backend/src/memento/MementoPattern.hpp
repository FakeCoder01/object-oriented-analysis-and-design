#pragma once
#include "../Shape.hpp"
#include <memory>
#include <stack>
#include <vector>

class Memento {
  std::vector<Shape> state;

public:
  Memento(const std::vector<Shape> &state);
  std::vector<Shape> getState() const;
};

class MementoEditor {
  std::vector<Shape> shapes;

public:
  void addShape(const Shape &s);
  void clear();
  std::unique_ptr<Memento> save();
  void restore(const Memento *m);
  json getShapes();
};

class Caretaker {
  std::stack<std::unique_ptr<Memento>> undo_stack;
  std::stack<std::unique_ptr<Memento>> redo_stack;
  MementoEditor *editor;

public:
  Caretaker(MementoEditor *ed);
  void backup();
  void undo();
  void redo();
};
