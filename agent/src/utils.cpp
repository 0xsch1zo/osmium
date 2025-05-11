#include "utils.hpp"
#include <array>
#include <stdexcept>
#include <windows.h>

std::string exec(const std::string &cmd) {
  auto pipe = popen((cmd + " 2>&1").c_str(), "r"); // get rid of shared_ptr

  if (!pipe)
    throw std::runtime_error("popen() failed!");

  std::string result;
  std::array<char, 128> buffer;
  while (!feof(pipe)) {
    if (fgets(buffer.data(), buffer.size(), pipe) != nullptr)
      result += buffer.data();
  }

  pclose(pipe);
  return result;
}
