cmake_minimum_required(VERSION 3.25)

set(CMAKE_EXPORT_COMPILE_COMMANDS ON)
set(CMAKE_CXX_STANDARD 23)
set(CMAKE_CXX_STANDARD_REQUIRED ON)
set(BUILD_SHARED_LIBS OFF)

include(FetchContent)

FetchContent_Declare(nullgate
    GIT_REPOSITORY https://github.com/0xsch1zo/NullGate
    GIT_TAG 1.1.2
)

FetchContent_Declare(cpr
    GIT_REPOSITORY https://github.com/libcpr/cpr
    GIT_TAG 1.11.2
)

FetchContent_MakeAvailable(cpr nullgate)

project(agent)

add_executable(agent
    ./src/agent.cpp
)

target_link_libraries(agent
    PRIVATE nullgate
    PRIVATE cpr::cpr
)
