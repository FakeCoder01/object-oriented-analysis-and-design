#pragma once
#include "json.hpp"
#include <string>

using json = nlohmann::json;

struct Shape {
  std::string type;
  int x, y;
  int size;
  std::string color;

  json toJson() const {
    return json{
        {"type", type}, {"x", x}, {"y", y}, {"size", size}, {"color", color}};
  }

  static Shape fromJson(const json &j) {
    return Shape{j["type"], j["x"], j["y"], j["size"], j["color"]};
  }
};
