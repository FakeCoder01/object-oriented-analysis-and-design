cmake_minimum_required(VERSION 3.10)
project(backend)

set(CMAKE_CXX_STANDARD 17)

find_package(Threads REQUIRED)

file(GLOB_RECURSE SOURCES "src/*.cpp")

add_executable(backend ${SOURCES})

# allow includes like "httplib.h" from the project root
target_include_directories(backend PRIVATE ${CMAKE_CURRENT_SOURCE_DIR})
target_link_libraries(backend Threads::Threads)
