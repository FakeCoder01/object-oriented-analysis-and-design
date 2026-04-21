#pragma once
#include "../Shape.hpp"
#include <stack>
#include <vector>

class NoPatternEditor {
  std::vector<Shape> current_shapes;
  std::stack<std::vector<Shape>> history;
  std::stack<std::vector<Shape>> redo_stack;

public:
  void addShape(const Shape &s);
  void undo();
  void redo();
  void clear();
  json getShapes();
};
